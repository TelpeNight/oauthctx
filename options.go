package oauthctx

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

type ConfigClientOptions struct {
	authClient    *http.Client
	requestClient *http.Client
}

func (o *ConfigClientOptions) TokenSourceOps() []TokenSourceOp {
	if o == nil || o.authClient == nil {
		return nil
	}
	return []TokenSourceOp{TokenSourceWithClient(o.authClient)}
}

func (o *ConfigClientOptions) ClientOps() []ClientOp {
	if o == nil || o.requestClient == nil {
		return nil
	}
	return []ClientOp{ClientWithRequestClient(o.requestClient)}
}

type ConfigClientOp func(o *ConfigClientOptions)

func BuildConfigClientOptions(ops ...ConfigClientOp) *ConfigClientOptions {
	if len(ops) == 0 {
		return nil
	}
	options := &ConfigClientOptions{}
	for _, op := range ops {
		op(options)
	}
	return options
}

func ConfigClientWithBaseClient(client *http.Client) ConfigClientOp {
	return func(o *ConfigClientOptions) {
		o.authClient = client
		o.requestClient = client
	}
}

func ConfigClientWithAuthClient(client *http.Client) ConfigClientOp {
	return func(o *ConfigClientOptions) {
		o.authClient = client
	}
}

func ConfigClientWithRequestClient(client *http.Client) ConfigClientOp {
	return func(o *ConfigClientOptions) {
		o.requestClient = client
	}
}

type TokenSourceOp func(o *TokenSourceOptions)

func TokenSourceWithClient(client *http.Client) TokenSourceOp {
	return func(o *TokenSourceOptions) {
		o.client = client
	}
}

// Oauth2ContextClient may be used with options. Returns nil on missing value
func Oauth2ContextClient(ctx context.Context) *http.Client {
	if ctx != nil {
		if hc, ok := ctx.Value(oauth2.HTTPClient).(*http.Client); ok {
			return hc
		}
	}
	return nil
}

type TokenSourceOptions struct {
	client *http.Client
}

func BuildTokenSourceOptions(ops ...TokenSourceOp) *TokenSourceOptions {
	if len(ops) == 0 {
		return nil
	}
	options := &TokenSourceOptions{}
	for _, op := range ops {
		op(options)
	}
	return options
}

func (o *TokenSourceOptions) WithOauth2HTTPClient(ctx context.Context) context.Context {
	if o == nil || o.client == nil {
		return ctx
	}
	return context.WithValue(ctx, oauth2.HTTPClient, o.client)
}

func (o *TokenSourceOptions) OpClient() *http.Client {
	if o == nil {
		return nil
	}
	return o.client
}
