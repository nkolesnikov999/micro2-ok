module github.com/nkolesnikov999/micro2-OK/iam

go 1.25.1

replace github.com/nkolesnikov999/micro2-OK/shared => ../shared

replace github.com/nkolesnikov999/micro2-OK/platform => ../platform

require (
	github.com/caarlos0/env/v11 v11.3.1
	github.com/gomodule/redigo v1.9.3
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.7.6
	github.com/joho/godotenv v1.5.1
	github.com/nkolesnikov999/micro2-OK/platform v0.0.0-00010101000000-000000000000
	github.com/nkolesnikov999/micro2-OK/shared v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.9.1
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.43.0
	google.golang.org/grpc v1.77.0
	google.golang.org/protobuf v1.36.10
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/pressly/goose/v3 v3.26.0 // indirect
	github.com/sethvargo/go-retry v0.3.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.46.1-0.20251013234738-63d1a5100f82 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8 // indirect
)
