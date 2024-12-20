package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/niradler/go-proxy-core/pkg/core"
	"github.com/niradler/go-ratelimit/pkg/rateLimiter"
)

type GithubShieldPlugin struct {
	allowedTokens  map[string]bool
	allowedMethods map[string]bool
	allowedIPs     map[string]bool
	rateLimiter    *rateLimiter.RateLimiter
}

func (p *GithubShieldPlugin) Init() error {
	p.allowedMethods = map[string]bool{
		"GET":    true,
		"POST":   true,
		"PUT":    true,
		"PATCH":  true,
		"DELETE": false, // Block DELETE requests
	}
	token := os.Getenv("GITHUB_ALLOWED_TOKEN")
	if token != "" {
		p.allowedTokens = map[string]bool{
			token: true,
		}
	}

	p.allowedIPs = map[string]bool{
		"::1": true,
	}

	return nil
}

func (p *GithubShieldPlugin) OnRequest(ctx *core.ProxyRequestContext) error {
	rctx := context.WithValue(ctx.Request.Context(), "request", ctx.Request)
	_, err := p.rateLimiter.Use(rctx)
	if err != nil {
		return err
	}

	authHeader := ctx.Request.Header.Get("Authorization")
	if authHeader == "" {
		return core.NewHttpError(fmt.Errorf("Authorization header missing"), 401)
	}
	if !p.allowedTokens[authHeader] {
		return core.NewHttpError(fmt.Errorf("Unauthorized request token"), 401)
	}

	method := ctx.Request.Method
	if !p.allowedMethods[method] {
		return core.NewHttpError(fmt.Errorf("Method %s not allowed", method), 403)
	}

	ip, _, err := net.SplitHostPort(ctx.Request.RemoteAddr)
	if err != nil {
		return core.NewHttpError(fmt.Errorf("Invalid IP address"), 400)
	}
	if !p.allowedIPs[ip] {
		return core.NewHttpError(fmt.Errorf("IP %s not allowed", ip), 403)
	}

	// Block requests that are not relevant to the vendor
	path := ctx.Request.URL.Path
	if !strings.HasPrefix(path, "/repos") && !strings.HasPrefix(path, "/issues") && !strings.HasPrefix(path, "/pulls") {
		return core.NewHttpError(fmt.Errorf("Path %s not allowed", path), 403)
	}
	return nil
}

func (a *GithubShieldPlugin) OnResponse(ctx *core.ProxyRequestContext) error {
	return nil
}

func (a *GithubShieldPlugin) OnError(ctx *core.ProxyRequestContext, err error) error {
	return err
}
