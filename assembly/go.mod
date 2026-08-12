module github.com/H1dEx/ms-rocket/assembly

go 1.26.4

require (
	github.com/H1dEx/ms-rocket/platform v0.0.0-00010101000000-000000000000
	github.com/H1dEx/ms-rocket/shared v0.0.0-00010101000000-000000000000
	github.com/IBM/sarama v1.60.1
	github.com/caarlos0/env/v11 v11.4.1
	github.com/go-faster/errors v0.8.0
	github.com/joho/godotenv v1.5.1
	go.uber.org/zap v1.28.0
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eapache/go-resiliency v1.7.0 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/klauspost/compress v1.19.1 // indirect
	github.com/pierrec/lz4/v4 v4.1.28 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20250401214520-65e299d6c5c9 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.54.0 // indirect
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260720211330-0afa2a65878a // indirect
	google.golang.org/grpc v1.82.1 // indirect
)

replace github.com/H1dEx/ms-rocket/shared => ../shared

replace github.com/H1dEx/ms-rocket/platform => ../platform
