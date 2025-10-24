package order

import (
	"context"
	"database/sql"
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
}

func (s *RepositorySuite) SetupSuite() {
	s.ctx = context.Background()

	// Сначала подключаемся к основной базе данных для создания тестовой БД
	mainURI := "postgres://order_user:order_password@localhost:5432/order_db"
	mainConn, err := pgx.Connect(s.ctx, mainURI)
	if err != nil {
		s.T().Fatalf("cannot connect to main PostgreSQL: %v", err)
	}
	defer func() {
		if closeErr := mainConn.Close(s.ctx); closeErr != nil {
			s.T().Logf("Failed to close main connection: %v", closeErr)
		}
	}()

	// Создаем тестовую базу данных
	s.testDBName = "order_test_" + time.Now().Format("20060102_150405")
	_, err = mainConn.Exec(s.ctx, "CREATE DATABASE "+s.testDBName)
	if err != nil {
		s.T().Fatalf("cannot create test database: %v", err)
	}

	// Подключаемся к тестовой базе данных
	testURI := "postgres://order_user:order_password@localhost:5432/" + s.testDBName
	conn, err := pgx.Connect(s.ctx, testURI)
	if err != nil {
		s.T().Fatalf("cannot connect to test PostgreSQL: %v", err)
	}
	s.conn = conn

	// Создаем sql.DB для миграций
	s.db = stdlib.OpenDB(*conn.Config())

	// Применяем миграции
	if err = goose.Up(s.db, "../../../migrations"); err != nil {
		s.T().Fatalf("cannot apply migrations: %v", err)
	}

	// Проверяем подключение
	if err = conn.Ping(s.ctx); err != nil {
		s.T().Fatalf("PostgreSQL ping failed: %v", err)
	}
}

func (s *RepositorySuite) TearDownSuite() {
	// Закрываем соединения
	if s.db != nil {
		if closeErr := s.db.Close(); closeErr != nil {
			s.T().Logf("Failed to close database connection: %v", closeErr)
		}
	}
	if s.conn != nil {
		if closeErr := s.conn.Close(s.ctx); closeErr != nil {
			s.T().Logf("Failed to close connection: %v", closeErr)
		}
	}

	// Удаляем тестовую базу данных
	if s.testDBName != "" {
		// Подключаемся к основной базе данных для удаления тестовой БД
		mainURI := "postgres://order_user:order_password@localhost:5432/order_db"
		mainConn, err := pgx.Connect(s.ctx, mainURI)
		if err != nil {
			s.T().Logf("Warning: cannot connect to main PostgreSQL for cleanup: %v", err)
			return
		}
		defer func() {
			if closeErr := mainConn.Close(s.ctx); closeErr != nil {
				s.T().Logf("Failed to close main connection in cleanup: %v", closeErr)
			}
		}()

		// Удаляем тестовую базу данных
		_, err = mainConn.Exec(s.ctx, "DROP DATABASE IF EXISTS "+s.testDBName)
		if err != nil {
			s.T().Logf("Warning: cannot drop test database %s: %v", s.testDBName, err)
		} else {
			s.T().Logf("Successfully dropped test database: %s", s.testDBName)
		}
	}
}

func (s *RepositorySuite) SetupTest() {
	// Очищаем таблицу orders перед каждым тестом
	if s.conn != nil {
		_, _ = s.conn.Exec(s.ctx, "DELETE FROM orders")
	}
	s.repository = NewRepository(s.conn)
}

func (s *RepositorySuite) TearDownTest() {
	// Очищаем таблицу orders после каждого теста
	if s.conn != nil {
		_, _ = s.conn.Exec(s.ctx, "DELETE FROM orders")
	}
}

func TestRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}
