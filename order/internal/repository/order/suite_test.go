package order

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"
)

type RepositorySuite struct {
	suite.Suite

	ctx        context.Context
	conn       *pgx.Conn
	db         *sql.DB
	repository *repository
	testDBName string
	pgURI      string
}

func (s *RepositorySuite) SetupSuite() {
	s.ctx = context.Background()

	// Базовый URI до БД (без имени базы для создания/удаления временных баз)
	s.pgURI = os.Getenv("POSTGRES_URI")
	if s.pgURI == "" {
		s.pgURI = "postgres://order_user:order_password@localhost:5432/order_db"
	}
}

func (s *RepositorySuite) TearDownSuite() {}

func (s *RepositorySuite) SetupTest() {
	// Создаем уникальную временную БД на каждый тест
	mainConn, err := pgx.Connect(s.ctx, s.pgURI)
	if err != nil {
		s.T().Fatalf("cannot connect to main PostgreSQL: %v", err)
	}
	defer func() { _ = mainConn.Close(s.ctx) }()

	s.testDBName = "order_test_" + time.Now().Format("20060102_150405_000000")
	if _, err := mainConn.Exec(s.ctx, "CREATE DATABASE "+s.testDBName); err != nil {
		s.T().Fatalf("cannot create test database: %v", err)
	}

	testURI := "postgres://order_user:order_password@localhost:5432/" + s.testDBName
	conn, err := pgx.Connect(s.ctx, testURI)
	if err != nil {
		s.T().Fatalf("cannot connect to test PostgreSQL: %v", err)
	}
	s.conn = conn
	s.db = stdlib.OpenDB(*conn.Config())

	if err := goose.Up(s.db, "../../../migrations"); err != nil {
		s.T().Fatalf("cannot apply migrations: %v", err)
	}
	if err := conn.Ping(s.ctx); err != nil {
		s.T().Fatalf("PostgreSQL ping failed: %v", err)
	}

	s.repository = NewRepository(s.conn)
}

func (s *RepositorySuite) TearDownTest() {
	// Закрываем соединения для текущей тестовой БД
	if s.db != nil {
		_ = s.db.Close()
		s.db = nil
	}
	if s.conn != nil {
		_ = s.conn.Close(s.ctx)
		s.conn = nil
	}

	// Удаляем текущую тестовую БД
	if s.testDBName != "" {
		mainConn, err := pgx.Connect(s.ctx, s.pgURI)
		if err == nil {
			_, _ = mainConn.Exec(s.ctx, "DROP DATABASE IF EXISTS "+s.testDBName)
			_ = mainConn.Close(s.ctx)
		}
		s.testDBName = ""
	}
}

func TestRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}
