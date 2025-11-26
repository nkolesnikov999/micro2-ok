package v1

import (
	"github.com/nkolesnikov999/micro2-OK/iam/internal/service"
	authV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/auth/v1"
)

type api struct {
	authV1.UnimplementedAuthServiceServer

	authService service.AuthService
}

func NewAPI(authService service.AuthService) *api {
	return &api{
		authService: authService,
	}
}
