module github.com/nkolesnikov999/micro2-OK/payment

go 1.25.1

replace github.com/nkolesnikov999/micro2-OK/shared => ../shared

replace github.com/nkolesnikov999/micro2-OK/platform => ../platform

require (
	github.com/brianvoe/gofakeit/v7 v7.8.0
	github.com/caarlos0/env/v11 v11.3.1
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/nkolesnikov999/micro2-OK/platform v0.0.0-00010101000000-000000000000
	github.com/nkolesnikov999/micro2-OK/shared v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.11.1
	google.golang.org/grpc v1.76.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
	google.golang.org/protobuf v1.36.10 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
