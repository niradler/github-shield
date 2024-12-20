module github.com/niadler/github-server-shield

go 1.23.3

require github.com/niradler/go-ratelimit v1.0.0

require github.com/niradler/go-proxy-core v1.0.0

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/redis/go-redis/v9 v9.7.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
)

replace github.com/niradler/go-ratelimit => C:\Projects\golang\ratelimit

replace github.com/niradler/go-proxy-core => C:\Projects\golang\go-proxy-core
