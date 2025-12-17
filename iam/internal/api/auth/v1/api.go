package v1

import (
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"

	"github.com/nkolesnikov999/micro2-OK/iam/internal/service"
	iamauth "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/auth/v1"
)

type api struct {
	iamauth.UnimplementedAuthServiceServer
	envoyauth.UnimplementedAuthorizationServer
	authService service.AuthService
}

func NewAPI(authService service.AuthService) *api {
	return &api{
		authService: authService,
	}
}
