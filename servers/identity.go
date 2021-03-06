package servers

import (
	"github.com/asaskevich/govalidator"
	"github.com/ellcrys/util"
	"github.com/ellcrys/patchain"
	"github.com/ellcrys/patchain/cockroach/tables"
	"github.com/ellcrys/patchain/object"
	"github.com/ellcrys/safehold/servers/common"
	"github.com/ellcrys/safehold/servers/proto_rpc"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

// getSystemIdentity returns system identity or a response error
func (s *RPC) getSystemIdentity() (*tables.Object, error) {
	obj, err := s.object.GetLast(&tables.Object{Key: object.MakeIdentityKey(systemEmail)})
	if err != nil {
		if err == patchain.ErrNotFound {
			return nil, common.NewSingleAPIErr(400, "", "", "system identity not found", nil)
		}
		return nil, common.ServerError
	}
	return obj, nil
}

// CreateIdentity creates a new identity for object
func (s *RPC) CreateIdentity(ctx context.Context, req *proto_rpc.CreateIdentityMsg) (*proto_rpc.GetObjectResponse, error) {

	// validate request
	if errs := validateIdentity(req); len(errs) > 0 {
		return nil, common.NewMultiAPIErr(400, "validation errors", errs)
	}

	systemIdentity, err := s.getSystemIdentity()
	if err != nil {
		logRPC.Errorf("%+v", err)
		return nil, err
	}

	objHandler := object.NewObject(s.db)

	// get existing identity
	existingIdentity, err := objHandler.GetLast(&tables.Object{OwnerID: systemIdentity.ID, Key: object.MakeIdentityKey(req.Email)})
	if err != nil {
		if err != patchain.ErrNotFound {
			return nil, common.ServerError
		}
	}

	// check if existing identity has been confirmed
	if err == nil {
		if existingIdentity.Ref3 == "confirmed" {
			return nil, common.NewSingleAPIErr(400, "used_email", "email", "Email is not available", nil)
		}
	}

	var numPartitions = int64(5)
	var newIdentity *tables.Object

	// create the identity and partitions
	var do = func() error {
		return s.db.Transact(func(txDB patchain.DB, commit patchain.CommitFunc, rollback patchain.RollbackFunc) error {

			var err error
			dbOpt := &patchain.UseDBOption{DB: txDB}
			passwordHash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

			if req.Developer {
				newIdentity = object.MakeDeveloperIdentityObject(systemIdentity.ID, systemIdentity.ID, req.Email, string(passwordHash), true)
			} else {
				newIdentity = object.MakeIdentityObject(systemIdentity.ID, systemIdentity.ID, req.Email, string(passwordHash), true)
			}

			err = s.object.Put(newIdentity, dbOpt)
			if err != nil {
				return err
			}

			_, err = s.object.CreatePartitions(numPartitions, newIdentity.ID, systemIdentity.ID, dbOpt)
			if err != nil {
				return err
			}

			return nil
		})
	}

	err = objHandler.Retry(func(stopRetry func()) error { return do() })
	if err != nil {
		logRPC.Errorf("%+v", err)
		return nil, common.ServerError
	}

	// TODO: send confirmation notification to email

	return &proto_rpc.GetObjectResponse{
		Object: util.MustStringify(newIdentity),
	}, nil
}

// GetIdentity fetches an identity object
func (s *RPC) GetIdentity(ctx context.Context, req *proto_rpc.GetIdentityMsg) (*proto_rpc.GetObjectResponse, error) {

	var err error
	var identity *tables.Object
	objHandler := object.NewObject(s.db)

	if govalidator.IsEmail(req.ID) {
		identity, err = objHandler.GetLast(&tables.Object{Key: object.MakeIdentityKey(req.ID)})
	} else {
		identity, err = objHandler.GetLast(&tables.Object{ID: req.ID, QueryParams: patchain.KeyStartsWith(object.IdentityPrefix)})
	}

	if err != nil {
		if err == patchain.ErrNotFound {
			return nil, common.NewSingleAPIErr(404, "", "", "Identity not found", nil)
		}
	}

	return &proto_rpc.GetObjectResponse{
		Object: util.MustStringify(identity),
	}, nil
}
