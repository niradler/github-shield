package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/niradler/go-proxy-core/pkg/core"
	"github.com/niradler/go-proxy-core/pkg/plugins"
	"github.com/niradler/go-ratelimit/pkg/rateLimiter"
	"github.com/niradler/go-ratelimit/pkg/store"
	"github.com/niradler/go-ratelimit/pkg/strategies"
)

func main() {
	// limit external vendors to 800 requests per hour per ip.
	usageRateLimiter := &rateLimiter.RateLimiter{
		Strategy: strategies.NewSlidingWindowStrategy(800, time.Hour),
		KeyFunc: func(ctx context.Context) (string, error) {
			r, ok := ctx.Value("request").(*http.Request)
			if !ok {
				return "", fmt.Errorf("could not get request from context")
			}

			return r.RemoteAddr, nil
		},
		DB: store.NewInMemoryStore(),
	}

	proxy, error := core.NewProxy()
	if error != nil {
		panic(error)
	}

	proxy.RegisterPlugin(&plugins.RequestIDPlugin{})
	proxy.RegisterPlugin(&plugins.LoggerPlugin{})

	err := proxy.RegisterPlugin(&GithubShieldPlugin{
		rateLimiter: usageRateLimiter,
		allowedTokens: map[string]bool{
			"secure-vendor-token-123": true,
		},
	})
	if err != nil {
		log.Fatalf("Failed to register GithubShieldPlugin: %v", err)
	}

	proxy.Start()
}
