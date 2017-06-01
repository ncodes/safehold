package servers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"net"

	"github.com/ellcrys/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"github.com/labstack/gommon/log"
	"github.com/ncodes/cocoon/core/config"
	"github.com/ncodes/patchain"
	"github.com/ncodes/safehold/servers/common"
	"github.com/ncodes/safehold/servers/oauth"
	"github.com/ncodes/safehold/servers/proto_rpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var logHTTP = config.MakeLogger("http")

// HTTP defines a structure for an HTTP server
// that provides REST API services.
type HTTP struct {
	rpcServerAddr string
	oauth         *oauth.OAuth
	db            patchain.DB
}

// NewHTTP creates an new http server instance
func NewHTTP(db patchain.DB, rpcServerAddr string) *HTTP {
	return &HTTP{
		rpcServerAddr: rpcServerAddr,
		db:            db,
		oauth:         oauth.NewOAuth(db),
	}
}

// getRouter returns the router
func (s *HTTP) getRouter() *mux.Router {
	r := mux.NewRouter()

	// oauth endpoints
	r.HandleFunc("/token", common.EasyHandle(http.MethodPost, s.oauth.GetToken))
	g := r.PathPrefix("/v1").Subrouter()

	// v1 endpoints
	g.HandleFunc("/identities", common.EasyHandle(http.MethodPost, s.createIdentity))
	g.HandleFunc("/sessions", common.EasyHandle(http.MethodPost, s.createSession))
	g.HandleFunc("/sessions/{id}", common.EasyHandle(http.MethodGet, s.getSession)).Methods(http.MethodGet)
	g.HandleFunc("/sessions/{id}", common.EasyHandle(http.MethodDelete, s.deleteSession)).Methods(http.MethodDelete)
	g.HandleFunc("/objects", common.EasyHandle(http.MethodPost, s.createObjects)).Methods(http.MethodPost)
	g.HandleFunc("/objects/query", common.EasyHandle(http.MethodPost, s.getObjects)).Methods(http.MethodPost)
	g.HandleFunc("/objects/count", common.EasyHandle(http.MethodPost, s.countObjects)).Methods(http.MethodPost)

	return r
}

// dialRPC creates a client connection to the RPC API
// and passes it to the callback function. It will close the connection
// after the callback function has finished
func (s *HTTP) dialRPC(cb func(client proto_rpc.APIClient) error) error {
	conn, err := grpc.Dial(s.rpcServerAddr, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "failed to dial RPC API")
	}
	defer conn.Close()
	return cb(proto_rpc.NewAPIClient(conn))
}

// Start starts the http server. Passes true to the startedCh channel
// when started
func (s *HTTP) Start(addr string, startedCB func(*HTTP)) error {

	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		logHTTP.Errorf("%+v", errors.Wrap(err, "failed to parse address"))
		return err
	}

	// attempt to connect to the RPC API
	if err := s.dialRPC(func(client proto_rpc.APIClient) error {
		return nil
	}); err != nil {
		logHTTP.Errorf("%+v", err)
		return err
	}

	time.AfterFunc(2*time.Second, func() {
		logHTTP.Infof("Started HTTP API server @ :%s", port)
		if startedCB != nil {
			startedCB(s)
		}
	})

	return http.ListenAndServe(addr, s.getRouter())
}

// createSession creates a new session
func (s *HTTP) createSession(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	var err error
	var body proto_rpc.DBSession
	var resp *proto_rpc.DBSession

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		return common.BodyMalformedError, 400
	}

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		var md metadata.MD
		authorization := r.Header.Get("Authorization")
		if len(authorization) > 0 {
			md = metadata.Pairs("authorization", authorization)
		}
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.CreateDBSession(ctx, &body)
		return err
	}); err != nil {
		logHTTP.Errorf("%+v", err)
		return err, 0
	}

	return resp, 201
}

// getSession gets a new session
func (s *HTTP) getSession(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	var err error
	var resp *proto_rpc.DBSession

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		var md metadata.MD
		authorization := r.Header.Get("Authorization")
		if len(authorization) > 0 {
			md = metadata.Pairs("authorization", authorization)
		}
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.GetDBSession(ctx, &proto_rpc.DBSession{
			ID: mux.Vars(r)["id"],
		})
		return err
	}); err != nil {
		logHTTP.Errorf("%+v", err)
		return err, 0
	}

	return resp, 201
}

// deleteSession deletes a session
func (s *HTTP) deleteSession(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	var err error
	var resp *proto_rpc.DBSession

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		var md metadata.MD
		authorization := r.Header.Get("Authorization")
		if len(authorization) > 0 {
			md = metadata.Pairs("authorization", authorization)
		}
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.DeleteDBSession(ctx, &proto_rpc.DBSession{
			ID: mux.Vars(r)["id"],
		})
		return err
	}); err != nil {
		logHTTP.Errorf("%+v", err)
		return err, 0
	}

	return resp, 201
}

// createIdentity creates an identity
func (s *HTTP) createIdentity(w http.ResponseWriter, r *http.Request) (interface{}, int) {

	var err error
	var body proto_rpc.CreateIdentityMsg
	var resp *proto_rpc.ObjectResponse

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		return common.BodyMalformedError, 400
	}

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		resp, err = client.CreateIdentity(context.Background(), &body)
		return err
	}); err != nil {
		log.Errorf("%+v", err)
		return err, 0
	}

	return resp, 201
}

// createObjects creates an arbitrary object.
// Requires an app token included in the authorization header
func (s *HTTP) createObjects(w http.ResponseWriter, r *http.Request) (interface{}, int) {

	var err error
	var body []*proto_rpc.Object
	var resp *proto_rpc.MultiObjectResponse

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		return common.BodyMalformedError, 400
	}

	var sessionID = r.URL.Query().Get("session")

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		var md = metadata.Join(metadata.Pairs("session_id", sessionID))
		authorization := r.Header.Get("Authorization")
		if len(authorization) > 0 {
			md = metadata.Join(md, metadata.Pairs("authorization", authorization))
		}
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Objects: body})
		return err
	}); err != nil {
		log.Errorf("%+v", err)
		return err, 0
	}

	return resp, 201
}

// getObjects performs query operations
func (s *HTTP) getObjects(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	var err error
	var resp *proto_rpc.MultiObjectResponse
	var rpcBody proto_rpc.GetObjectMsg
	var body struct {
		Query   map[string]interface{} `json:"query"`
		Owner   string                 `json:"owner"`
		Creator string                 `json:"creator"`
		Limit   int32                  `json:"limit"`
		Order   []struct {
			Field string `json:"field"`
			Order int32  `json:"order"`
		} `json:"order"`
	}

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		return common.BodyMalformedError, 400
	}

	copier.Copy(&rpcBody, body)
	rpcBody.Query = util.MustStringify(body.Query)

	var sessionID = r.URL.Query().Get("session")

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		var md = metadata.Join(metadata.Pairs("session_id", sessionID))
		authorization := r.Header.Get("Authorization")
		if len(authorization) > 0 {
			md = metadata.Join(md, metadata.Pairs("authorization", authorization))
		}
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.GetObjects(ctx, &rpcBody)
		return err
	}); err != nil {
		log.Errorf("%+v", err)
		return err, 0
	}

	return resp, 200
}

// countObjects counts objects matching a query
func (s *HTTP) countObjects(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	var err error
	var resp *proto_rpc.ObjectCountResponse
	var rpcBody proto_rpc.GetObjectMsg
	var body struct {
		Query   map[string]interface{} `json:"query"`
		Owner   string                 `json:"owner"`
		Creator string                 `json:"creator"`
	}

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		return common.BodyMalformedError, 400
	}

	copier.Copy(&rpcBody, body)
	rpcBody.Query = util.MustStringify(body.Query)

	var sessionID = r.URL.Query().Get("session")

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		var md = metadata.Join(metadata.Pairs("session_id", sessionID))
		authorization := r.Header.Get("Authorization")
		if len(authorization) > 0 {
			md = metadata.Join(md, metadata.Pairs("authorization", authorization))
		}
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.CountObjects(ctx, &rpcBody)
		return err
	}); err != nil {
		log.Errorf("%+v", err)
		return err, 0
	}

	return resp, 200
}
