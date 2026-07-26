module github.com/H1dEx/ms-rocket/inventory

go 1.26.4

require (
	github.com/H1dEx/ms-rocket/platform v0.0.0-00010101000000-000000000000
	github.com/H1dEx/ms-rocket/shared v0.0.0-20260725113859-5e4cb1b848d1
	github.com/caarlos0/env/v11 v11.4.1
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/samber/lo v1.53.0
	github.com/stretchr/testify v1.11.1
	go.mongodb.org/mongo-driver v1.17.9
	go.uber.org/zap v1.28.0
	google.golang.org/grpc v1.82.1
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/golang/snappy v1.0.0 // indirect
	github.com/klauspost/compress v1.19.1 // indirect
	github.com/montanaflynn/stats v0.12.2 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.2.0 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.54.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
)

require (
	github.com/brianvoe/gofakeit/v7 v7.15.0
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260724162435-b2f20204f0df // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/H1dEx/ms-rocket/shared => ../shared

replace github.com/H1dEx/ms-rocket/platform => ../platform
