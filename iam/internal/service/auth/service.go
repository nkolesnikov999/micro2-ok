package auth

import (
	"github.com/nkolesnikov999/micro2-OK/iam/internal/repository"
	def "github.com/nkolesnikov999/micro2-OK/iam/internal/service"
)

var _ def.AuthService = (*service)(nil)

type service struct {
	sessionRepository repository.SessionRepository
	userRepository    repository.UserRepository
}

func NewService(
	sessionRepository repository.SessionRepository,
	userRepository repository.UserRepository,
) *service {
	return &service{
		sessionRepository: sessionRepository,
		userRepository:    userRepository,
	}
}
