package user

import (
	"github.com/jackc/pgx/v5"

	def "github.com/nkolesnikov999/micro2-OK/iam/internal/repository"
)

var _ def.UserRepository = (*repository)(nil)

type repository struct {
	connDB *pgx.Conn
}

func NewRepository(connDB *pgx.Conn) *repository {
	return &repository{
		connDB: connDB,
	}
}
