package main

import (
	"context"
	"fmt"

	"github.com/a2aproject/a2a-go/a2aclient"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

type AuthInterceptor struct {
	tokenSource oauth2.TokenSource
}

func (a AuthInterceptor) Before (ctx context.Context, req *a2aclient.Request) (context.Context, error) {

	// Create an auth token
	authToken, err := a.tokenSource.Token()
	if err != nil {
		return ctx, fmt.Errorf("tokenSource.Token: %w", err)
	}

	// Set the authorization header
	if req.Meta == nil {
        req.Meta = make(map[string][]string)
    }
    req.Meta["Authorization"]= []string{"Bearer " + authToken.AccessToken}

	return ctx, nil 
}

func (a AuthInterceptor) After (ctx context.Context, resp *a2aclient.Response) error {
	return nil
}

// NewAuthInterceptor creates a new request interceptor satisfying hte [a2aclient.CallInterceptor] interface
func NewAuthInterceptor(audience string) (AuthInterceptor, error) {
	ctx := context.Background()

	// Fetch environment defalt credentials
	credentials, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return AuthInterceptor{}, fmt.Errorf("failed to generate default credentials: %w", err)
	}

	// Establish a token source
	ts, err := idtoken.NewTokenSource(ctx, audience, option.WithCredentials(credentials))
	if err != nil {
		return AuthInterceptor{}, fmt.Errorf("idtoken.NewTokenSource: %w", err)
	}

	return AuthInterceptor{
		tokenSource: ts,
	}, nil
}