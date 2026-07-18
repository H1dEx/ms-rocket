module github.com/H1dEx/ms-rocket/inventory

go 1.25.0

require (
	github.com/H1dEx/ms-rocket/shared v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.11.1
	google.golang.org/grpc v1.82.0
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/brianvoe/gofakeit/v7 v7.15.0
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260713224248-f5fc221cf8c4 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/H1dEx/ms-rocket/shared => ../shared
