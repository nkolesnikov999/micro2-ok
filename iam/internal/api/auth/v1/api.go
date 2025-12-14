package v1

import (
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/nkolesnikov999/micro2-OK/iam/internal/service"
	authV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/auth/v1"
)

type api struct {
	authV1.UnimplementedAuthServiceServer
	authv3.UnimplementedAuthorizationServer
	authService service.AuthService
}

func NewAPI(authService service.AuthService) *api {
	return &api{
		authService: authService,
	}
}
