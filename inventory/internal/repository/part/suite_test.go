package part

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RepositorySuite struct {
	suite.Suite

	ctx        context.Context
	client     *mongo.Client
	db         *mongo.Database
	repository *repository
}

func (s *RepositorySuite) SetupSuite() {
	s.ctx = context.Background()

	mongoURI := "mongodb://inventory_user:inventory_password@localhost:27017"

	client, err := mongo.Connect(s.ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		s.T().Fatalf("cannot connect to Mongo: %v", err)
	}
	s.client = client

	if err = client.Ping(s.ctx, nil); err != nil {
		s.T().Fatalf("Mongo ping failed: %v", err)
	}

	testDBName := "inventory_test_" + time.Now().Format("20060102_150405")
	s.db = client.Database(testDBName)
}

func (s *RepositorySuite) TearDownSuite() {
	if s.client != nil && s.db != nil {
		_ = s.db.Drop(s.ctx)
		_ = s.client.Disconnect(s.ctx)
	}
}

func (s *RepositorySuite) SetupTest() {
	s.repository = NewRepository(s.ctx, s.db)
}

func (s *RepositorySuite) TearDownTest() {
	if s.db != nil {
		_ = s.db.Collection("parts").Drop(s.ctx)
	}
}

func TestRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}
