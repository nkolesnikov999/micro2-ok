package v1

import (
	"github.com/nkolesnikov999/micro2-OK/iam/internal/service"
	userV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/user/v1"
)

type api struct {
	userV1.UnimplementedUserServiceServer

	userService service.UserService
}

func NewAPI(userService service.UserService) *api {
	return &api{
		userService: userService,
	}
}
