module github.com/H1dEx/ms-rocket/payment

go 1.25.0

require (
	github.com/H1dEx/ms-rocket/shared v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.82.0
)

require (
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260713224248-f5fc221cf8c4 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/H1dEx/ms-rocket/shared => ../shared
