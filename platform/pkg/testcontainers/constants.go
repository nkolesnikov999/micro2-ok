package testcontainers

// MongoDB constants
const (
	// MongoDB container constants
	MongoContainerName = "mongo"
	MongoPort          = "27017"

	// MongoDB environment variables
	MongoImageNameKey = "MONGO_IMAGE_NAME"
	MongoHostKey      = "MONGO_HOST"
	MongoPortKey      = "MONGO_PORT"
	MongoDatabaseKey  = "MONGO_DATABASE"
	MongoUsernameKey  = "MONGO_INITDB_ROOT_USERNAME"
	MongoPasswordKey  = "MONGO_INITDB_ROOT_PASSWORD" //nolint:gosec
	MongoAuthDBKey    = "MONGO_AUTH_DB"
)

// PostgreSQL constants
const (
	// PostgreSQL container constants
	PostgresContainerName = "postgres"
	PostgresPort          = "5432"

	// PostgreSQL environment variables
	PostgresImageNameKey = "POSTGRES_IMAGE_NAME"
	PostgresHostKey      = "POSTGRES_HOST"
	PostgresPortKey      = "POSTGRES_PORT"
	PostgresDatabaseKey  = "POSTGRES_DB"
	PostgresUserKey      = "POSTGRES_USER"
	PostgresPasswordKey  = "POSTGRES_PASSWORD" //nolint:gosec
	PostgresSSLModeKey   = "POSTGRES_SSL_MODE"
)

// Redis constants
const (
	// Redis container constants
	RedisContainerName = "redis"
	RedisPort          = "6379"

	// Redis environment variables
	RedisImageNameKey = "REDIS_IMAGE_NAME"
	RedisHostKey      = "REDIS_HOST"
	RedisPortKey      = "REDIS_PORT"
)
