package servers

import (
	"fmt"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"net"

	"github.com/ellcrys/util"
	"github.com/ellcrys/safehold/servers/common"
	"github.com/ellcrys/safehold/servers/proto_rpc"
	"github.com/ellcrys/safehold/session"
	"golang.org/x/net/context"
)

// CreateDBSession creates a new session and returns the session ID
func (s *RPC) CreateDBSession(ctx context.Context, req *proto_rpc.DBSession) (*proto_rpc.DBSession, error) {

	developerID := ctx.Value(CtxIdentity).(string)

	sessionID := util.UUID4()
	sessionID, err := s.dbSession.CreateSession(sessionID, developerID)
	if err != nil {
		logRPC.Error("%+v", err)
		return nil, common.NewSingleAPIErr(500, "", "", "session not created", nil)
	}

	return &proto_rpc.DBSession{
		ID: sessionID,
	}, nil
}

// GetDBSession gets a database session.
func (s *RPC) GetDBSession(ctx context.Context, req *proto_rpc.DBSession) (*proto_rpc.DBSession, error) {
	developerID := ctx.Value(CtxIdentity).(string)

	if req.ID == "" {
		return nil, common.NewSingleAPIErr(400, "", "", "session id is required", nil)
	}

	sessionID := req.ID

	// check if session exists locally
	if sessionAgent := s.dbSession.GetAgent(sessionID); sessionAgent != nil {

		// check if session is owned by the developer, if not, return permission error
		if sessionAgent.OwnerID != developerID {
			return nil, common.NewSingleAPIErr(401, "", "", "permission denied: you don't have permission to perform this operation", nil)
		}

		return req, nil
	}

	// check session registry
	regItem, err := s.sessionReg.Get(sessionID)
	if err != nil {
		if err == session.ErrNotFound {
			return nil, common.NewSingleAPIErr(404, "", "", "session not found", nil)
		}
	}

	// check if session is owned by the developer, if not, return permission error
	if regItem.Meta["identity"] != developerID {
		return nil, common.NewSingleAPIErr(401, "", "", "permission denied: you don't have permission to perform this operation", nil)
	}

	return req, nil
}

// DeleteDBSession deletes a existing database session
func (s *RPC) DeleteDBSession(ctx context.Context, req *proto_rpc.DBSession) (*proto_rpc.DBSession, error) {

	developerID := ctx.Value(CtxIdentity).(string)
	authorization := util.FromIncomingMD(ctx, "authorization")
	localOnly := util.FromIncomingMD(ctx, "local-only") == "true"

	if req.ID == "" {
		return nil, common.NewSingleAPIErr(400, "", "", "session id is required", nil)
	}

	sessionID := req.ID

	// check if session exists locally, if so, delete immediately
	if sessionAgent := s.dbSession.GetAgent(sessionID); sessionAgent != nil {

		// check if session is owned by the developer, if not, return permission error
		if sessionAgent.OwnerID != developerID {
			return nil, common.NewSingleAPIErr(401, "", "", "permission denied: you don't have permission to perform this operation", nil)
		}

		s.dbSession.End(sessionID)
		return req, nil
	}

	// if localOnly is true, return error
	if localOnly {
		return nil, fmt.Errorf("session not found")
	}

	// get session from the session registry.
	// if we find it, we need to call the session's host server to delete.
	// if not found, return `not found` error
	item, err := s.sessionReg.Get(sessionID)
	if err != nil {
		if err == session.ErrNotFound {
			return nil, common.NewSingleAPIErr(404, "", "", "session not found", nil)
		}
		return nil, common.ServerError
	}

	// check if session is owned by the developer, if not, return permission error
	if item.Meta["identity"] != developerID {
		return nil, common.NewSingleAPIErr(401, "", "", "permission denied: you don't have permission to perform this operation", nil)
	}

	sessionHostAddr := net.JoinHostPort(item.Address, strconv.Itoa(item.Port))
	client, err := grpc.Dial(sessionHostAddr, grpc.WithInsecure())
	if err != nil {
		return nil, common.ServerError
	}
	defer client.Close()

	// make call to the session host server.
	// include the auth token of the current request
	// and set local-only to force the RPC method
	// to only check local/in-memory session cache
	server := proto_rpc.NewAPIClient(client)
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", authorization, "local-only", "true"))
	resp, err := server.DeleteDBSession(ctx, req)
	if err != nil {
		if grpc.ErrorDesc(err) == "session not found" {
			return nil, common.NewSingleAPIErr(404, "", "", "session not found", nil)
		}
		logRPC.Errorf("%+v", err)
		return nil, common.ServerError
	}

	return resp, nil
}

// CommitSession commits an active session. Returns no error if the session was successfully
// committed, or has already been committed or does not exist.
func (s *RPC) CommitSession(ctx context.Context, req *proto_rpc.DBSession) (*proto_rpc.DBSession, error) {

	developerID := ctx.Value(CtxIdentity).(string)
	authorization := util.FromIncomingMD(ctx, "authorization")
	localOnly := util.FromIncomingMD(ctx, "local-only") == "true"

	if req.ID == "" {
		return nil, common.NewSingleAPIErr(400, "", "", "session id is required", nil)
	}

	sessionID := req.ID

	// check if session exists locally, if so, authenticate developer and commit immediately
	if sessionAgent := s.dbSession.GetAgent(sessionID); sessionAgent != nil {

		// check if session is owned by the developer, if not, return permission error
		if sessionAgent.OwnerID != developerID {
			return nil, common.NewSingleAPIErr(401, "", "", "permission denied: you don't have permission to perform this operation", nil)
		}

		s.dbSession.Commit(sessionID)
		return req, nil
	}

	// if localOnly is true, return error
	if localOnly {
		return nil, fmt.Errorf("session not found")
	}

	// get session from the session registry.
	// if we find it, we need to call the session's host server to delete.
	// if not found, return `not found` error
	item, err := s.sessionReg.Get(sessionID)
	if err != nil {
		if err == session.ErrNotFound {
			return nil, common.NewSingleAPIErr(404, "", "", "session not found", nil)
		}
		return nil, common.ServerError
	}

	// check if session is owned by the developer, if not, return permission error
	if item.Meta["identity"] != developerID {
		return nil, common.NewSingleAPIErr(401, "", "", "permission denied: you don't have permission to perform this operation", nil)
	}

	sessionHostAddr := net.JoinHostPort(item.Address, strconv.Itoa(item.Port))
	client, err := grpc.Dial(sessionHostAddr, grpc.WithInsecure())
	if err != nil {
		return nil, common.ServerError
	}
	defer client.Close()

	// make call to the session host server.
	// include the auth token of the current request
	// and set local-only to force the RPC method
	// to only check local/in-memory session cache
	server := proto_rpc.NewAPIClient(client)
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", authorization, "local-only", "true"))
	resp, err := server.CommitSession(ctx, req)
	if err != nil {
		if grpc.ErrorDesc(err) == "session not found" {
			return nil, common.NewSingleAPIErr(404, "", "", "session not found", nil)
		}
		logRPC.Errorf("%+v", err)
		return nil, common.ServerError
	}

	return resp, nil
}

// RollbackSession rolls back a session. Returns no error if the session was successfully
// rolled back, or has already been rolled back or does not exist.
func (s *RPC) RollbackSession(ctx context.Context, req *proto_rpc.DBSession) (*proto_rpc.DBSession, error) {

	developerID := ctx.Value(CtxIdentity).(string)
	authorization := util.FromIncomingMD(ctx, "authorization")
	localOnly := util.FromIncomingMD(ctx, "local-only") == "true"

	if req.ID == "" {
		return nil, common.NewSingleAPIErr(400, "", "", "session id is required", nil)
	}

	sessionID := req.ID

	// check if session exists locally, if so, authenticate developer and rollback immediately
	if sessionAgent := s.dbSession.GetAgent(sessionID); sessionAgent != nil {

		// check if session is owned by the developer, if not, return permission error
		if sessionAgent.OwnerID != developerID {
			return nil, common.NewSingleAPIErr(401, "", "", "permission denied: you don't have permission to perform this operation", nil)
		}

		s.dbSession.Rollback(sessionID)
		return req, nil
	}

	// if localOnly is true, return error
	if localOnly {
		return nil, fmt.Errorf("session not found")
	}

	// get session from the session registry.
	// if we find it, we need to call the session's host server to delete.
	// if not found, return `not found` error
	item, err := s.sessionReg.Get(sessionID)
	if err != nil {
		if err == session.ErrNotFound {
			return nil, common.NewSingleAPIErr(404, "", "", "session not found", nil)
		}
		return nil, common.ServerError
	}

	// check if session is owned by the developer, if not, return permission error
	if item.Meta["identity"] != developerID {
		return nil, common.NewSingleAPIErr(401, "", "", "permission denied: you don't have permission to perform this operation", nil)
	}

	sessionHostAddr := net.JoinHostPort(item.Address, strconv.Itoa(item.Port))
	client, err := grpc.Dial(sessionHostAddr, grpc.WithInsecure())
	if err != nil {
		return nil, common.ServerError
	}
	defer client.Close()

	// make call to the session host server.
	// include the auth token of the current request
	// and set local-only to force the RPC method
	// to only check local/in-memory session cache
	server := proto_rpc.NewAPIClient(client)
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", authorization, "local-only", "true"))
	resp, err := server.RollbackSession(ctx, req)
	if err != nil {
		if grpc.ErrorDesc(err) == "session not found" {
			return nil, common.NewSingleAPIErr(404, "", "", "session not found", nil)
		}
		logRPC.Errorf("%+v", err)
		return nil, common.ServerError
	}

	return resp, nil
}
