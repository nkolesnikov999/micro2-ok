module github.com/nkolesnikov999/micro2-OK/inventory

go 1.25.1

replace github.com/nkolesnikov999/micro2-OK/shared => ../shared

require (
	github.com/brianvoe/gofakeit/v7 v7.7.3
	github.com/google/uuid v1.6.0
	github.com/nkolesnikov999/micro2-OK/shared v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.76.0
	google.golang.org/protobuf v1.36.10
)

require (
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
)
