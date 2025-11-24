package user

import (
	"github.com/nkolesnikov999/micro2-OK/iam/internal/repository"
	def "github.com/nkolesnikov999/micro2-OK/iam/internal/service"
)

var _ def.UserService = (*service)(nil)

type service struct {
	userRepository repository.UserRepository
}

func NewService(
	userRepository repository.UserRepository,
) *service {
	return &service{
		userRepository: userRepository,
	}
}
