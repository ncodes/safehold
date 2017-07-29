package servers

import (
	"fmt"

	"google.golang.org/grpc"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ellcrys/elldb/servers/common"
	"github.com/ellcrys/elldb/servers/db"
	"github.com/ellcrys/elldb/servers/oauth"
	"github.com/ellcrys/elldb/types"
	"github.com/ellcrys/util"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jinzhu/gorm"
	"golang.org/x/net/context"
)

var (
	// CtxAccount represents an authenticated account
	CtxAccount types.CtxKey = "account"

	// CtxTokenClaims represents claims in an auth token
	CtxTokenClaims types.CtxKey = "token_claims"

	// ErrInvalidToken represents a error about an invalid token
	ErrInvalidToken = fmt.Errorf("permission denied. Invalid token")

	// methodsNotRequiringAuth includes the full method name of methods that
	// must not be processed by the auth interceptor
	methodsNotRequiringAuth = []string{
		"/proto_rpc.API/CreateAccount",
	}

	// methodsRequiringAppToken includes the full method name of methods that
	// must pass app token authorization checks
	methodsRequiringAppToken = []string{
		"/proto_rpc.API/GetAccount",
		"/proto_rpc.API/CreateBucket",
		"/proto_rpc.API/CreateObjects",
		"/proto_rpc.API/GetObjects",
		"/proto_rpc.API/CreateSession",
		"/proto_rpc.API/GetSession",
		"/proto_rpc.API/DeleteSession",
		"/proto_rpc.API/CountObjects",
		"/proto_rpc.API/CreateMapping",
		"/proto_rpc.API/GetMapping",
		"/proto_rpc.API/GetAllMapping",
		"/proto_rpc.API/CommitSession",
		"/proto_rpc.API/RollbackSession",
		"/proto_rpc.API/UpdateObjects",
		"/proto_rpc.API/DeleteObjects",
	}
)

// Interceptors returns the API interceptors
func (s *RPC) Interceptors() grpc.UnaryServerInterceptor {
	return middleware.ChainUnaryServer(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logRPC.Debugf("New request [method=%s]", info.FullMethod)
		return handler(ctx, req)
	}, s.authInterceptor, s.requiresAppTokenInterceptor)
}

// authInterceptor checks whether the request has valid access token
func (s *RPC) authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	if util.InStringSlice(methodsNotRequiringAuth, info.FullMethod) {
		return handler(ctx, req)
	}

	tokenStr, err := util.GetAuthToken(ctx, "bearer")
	if err != nil {
		return nil, common.NewSingleAPIErr(401, common.CodeAuthorizationError, "", err.Error(), nil)
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		tokenType := token.Claims.(jwt.MapClaims)["type"]
		if tokenType == oauth.TokenTypeApp {
			return []byte(oauth.SigningSecret), nil
		}

		return nil, fmt.Errorf("unknown token type")
	})
	if err != nil {
		logRPC.Debugf("%+v", err)
		return nil, common.NewSingleAPIErr(401, common.CodeAuthorizationError, "", ErrInvalidToken.Error(), nil)
	}

	claims := token.Claims.(jwt.MapClaims)
	if err = claims.Valid(); err != nil {
		return nil, common.NewSingleAPIErr(401, common.CodeAuthorizationError, "", ErrInvalidToken.Error(), nil)
	}

	ctx = context.WithValue(ctx, CtxTokenClaims, claims)

	return handler(ctx, req)
}

// requiresAppTokenInterceptor checks if the claim produced by authInterceptor
// is a valid app token. This will only apply to methods that require app token. It immediately
// calls the next handler if the method does not require app token.
func (s *RPC) requiresAppTokenInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	var err error

	// verify claims for methods that require an app token
	if !util.InStringSlice(methodsRequiringAppToken, info.FullMethod) {
		return handler(ctx, req)
	}

	claims := ctx.Value(CtxTokenClaims).(jwt.MapClaims)
	if claims["type"] != oauth.TokenTypeApp {
		return nil, common.NewSingleAPIErr(401, "", "", "endpoint requires an app token", nil)
	}

	ctx, err = s.processAppTokenClaims(ctx, claims)
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

// processAppTokenClaims checks whether a token claim is valid. It confirms the
// the account associated with the claim and includes the account in the context.
func (s *RPC) processAppTokenClaims(ctx context.Context, claims jwt.MapClaims) (context.Context, error) {

	// get account with matching client id
	var account db.Account
	if err := s.db.Where("client_id = ?", claims["id"].(string)).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx, common.NewSingleAPIErr(401, common.CodeAuthorizationError, "", ErrInvalidToken.Error(), nil)
		}
		logRPC.Errorf("%+v", err)
		return nil, common.ServerError
	}

	return context.WithValue(ctx, CtxAccount, account.ID), nil
}
