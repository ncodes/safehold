package servers

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/ellcrys/util"
	"github.com/labstack/gommon/log"
	"github.com/ellcrys/patchain"
	"github.com/ellcrys/patchain/cockroach"
	"github.com/ellcrys/patchain/cockroach/tables"
	"github.com/ellcrys/patchain/object"
	"github.com/ellcrys/safehold/servers/common"
	"github.com/ellcrys/safehold/servers/proto_rpc"
	"github.com/ellcrys/safehold/session"
	"golang.org/x/net/context"
)

// validateMapping checks whether a map is valid
func validateMapping(mapping map[string]string) []common.Error {
	var errs []common.Error
	validObjectFields := cockroach.NewDB().GetValidObjectFields()

	if len(mapping) == 0 {
		errs = append(errs, common.Error{
			Field:   "mapping",
			Message: "requires at least one mapped field",
		})
	}

	for mapField, fieldData := range mapping {
		if len(fieldData) == 0 {
			errs = append(errs, common.Error{
				Field:   mapField,
				Message: fmt.Sprintf("column name is required"),
			})
			continue
		}
		if !util.InStringSlice(validObjectFields, fieldData) {
			errs = append(errs, common.Error{
				Field:   mapField,
				Message: fmt.Sprintf("column name '%s' is unknown", fieldData),
			})
		}
	}
	return errs
}

// getMappingWithSession gets a mapping using an active session. It will
// find mapping belonging to the ownerID if the mapping name is not a UUIDv4
// string, otherwise it will find any mapping matching the mapping name.
func (s *RPC) getMappingWithSession(sid, mappingName, ownerID string) (map[string]string, error) {

	var query = tables.Object{Key: object.MakeMappingKey(mappingName), OwnerID: ownerID}
	if govalidator.IsUUIDv4(mappingName) {
		query = tables.Object{ID: mappingName}
	}

	var mappingObj tables.Object
	if err := session.SendQueryOpWithSession(s.dbSession, sid, "", &query, 1, "", &mappingObj); err != nil {
		if err == patchain.ErrNotFound {
			return nil, common.NewSingleAPIErr(404, "", "", "mapping with name=`"+mappingName+"` does not exist", nil)
		}
		return nil, common.ServerError
	}

	var mapping map[string]string
	if err := util.FromJSON([]byte(mappingObj.Value), &mapping); err != nil {
		return nil, common.NewSingleAPIErr(500, "", "", "failed to parse mapping value", nil)
	}

	return mapping, nil
}

// CreateMapping creates a mapping for an identity
func (s *RPC) CreateMapping(ctx context.Context, req *proto_rpc.CreateMappingMsg) (*proto_rpc.CreateMappingResponse, error) {

	developerID := ctx.Value(CtxIdentity).(string)

	if len(req.Name) == 0 {
		req.Name = util.UUID4()
	}

	var mapping map[string]string
	if err := util.FromJSON(req.Mapping, &mapping); err != nil {
		return nil, common.NewSingleAPIErr(404, "", "", "malformed mapping. Expected json object.", nil)
	}

	// validate mapping
	errs := validateMapping(mapping)
	if len(errs) > 0 {
		return nil, common.NewMultiAPIErr(400, "mapping error", errs)
	}

	sessionID := util.UUID4()
	s.dbSession.CreateUnregisteredSession(sessionID, developerID)
	defer s.dbSession.CommitEnd(sessionID)

	mappingObj := object.MakeMappingObject(developerID, req.Name, string(req.Mapping))
	err := s.dbSession.SendOp(sessionID, &session.Op{
		OpType:  session.OpPutObjects,
		PutData: mappingObj,
		Done:    make(chan struct{}),
	})
	if err != nil {
		s.dbSession.RollbackEnd(sessionID)
		return nil, common.ServerError
	}

	return &proto_rpc.CreateMappingResponse{
		Name: req.Name,
		ID:   mappingObj.ID,
	}, nil
}

// GetMapping fetches a mapping belonging to the logged in developer
func (s *RPC) GetMapping(ctx context.Context, req *proto_rpc.GetMappingMsg) (*proto_rpc.GetMappingResponse, error) {

	developerID := ctx.Value(CtxIdentity).(string)

	if len(req.Name) == 0 {
		return nil, common.NewSingleAPIErr(400, "", "id", "name is required", nil)
	}

	sessionID := util.UUID4()
	s.dbSession.CreateUnregisteredSession(sessionID, developerID)
	defer s.dbSession.CommitEnd(sessionID)

	// construct query. Find mapping with matching name/key that belongs to the current developer
	// if the request name is not a UUID. Otherwise find any mapping with the matching UUID owned by any developer.
	query := tables.Object{Key: object.MakeMappingKey(req.Name), OwnerID: developerID}
	if govalidator.IsUUIDv4(req.Name) {
		query.Key = ""
		query.ID = req.Name
		query.OwnerID = ""
		query.QueryParams = patchain.KeyStartsWith(object.MappingPrefix)
	}

	var mapping tables.Object
	err := s.dbSession.SendOp(sessionID, &session.Op{
		OpType:          session.OpGetObjects,
		QueryWithObject: &query,
		OrderBy:         "timestamp desc",
		Limit:           1,
		Out:             &mapping,
		Done:            make(chan struct{}),
	})
	if err != nil {
		if err == patchain.ErrNotFound {
			return nil, common.NewSingleAPIErr(404, "", "", "mapping not found", nil)
		}
		log.Errorf("%+v", err)
		return nil, common.ServerError
	}

	return &proto_rpc.GetMappingResponse{
		Mapping: util.MustStringify(mapping),
	}, nil
}

// GetAllMapping fetches the most recent mappings belonging to the logged in developer
func (s *RPC) GetAllMapping(ctx context.Context, req *proto_rpc.GetAllMappingMsg) (*proto_rpc.GetAllMappingResponse, error) {
	developerID := ctx.Value(CtxIdentity).(string)

	if req.Limit == 0 {
		req.Limit = 50
	}

	sessionID := util.UUID4()
	s.dbSession.CreateUnregisteredSession(sessionID, developerID)
	defer s.dbSession.CommitEnd(sessionID)

	// create query. Find all mapping belonging to the developer. If mapping name is set, find
	// mapping with the matching name that belongs to the developer.
	var query = tables.Object{QueryParams: patchain.KeyStartsWith(object.MappingPrefix), OwnerID: developerID}
	if len(req.Name) > 0 {
		query.QueryParams = patchain.QueryParams{}
		query.Key = object.MakeMappingKey(req.Name)
	}

	var mappings []tables.Object
	err := s.dbSession.SendOp(sessionID, &session.Op{
		OpType:          session.OpGetObjects,
		QueryWithObject: &query,
		OrderBy:         "timestamp desc",
		Limit:           int(req.Limit),
		Out:             &mappings,
		Done:            make(chan struct{}),
	})
	if err != nil {
		if err == patchain.ErrNotFound {
			return nil, common.NewSingleAPIErr(404, "", "", "mapping not found", nil)
		}
		log.Errorf("%+v", err)
		return nil, common.ServerError
	}

	return &proto_rpc.GetAllMappingResponse{
		Mappings: util.MustStringify(mappings),
	}, nil
}
