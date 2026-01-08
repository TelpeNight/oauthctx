package oauthctx

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

// RequestFlowConfig is configuration for oauth2 request flow.
//
// Currently, provides optional http clients
type RequestFlowConfig struct {
	tokenClient   *http.Client
	requestClient *http.Client
}

// TokenSourceOps creates options for functions like Config.TokenSource
func (o *RequestFlowConfig) TokenSourceOps() []TokenSourceOp {
	if o == nil || o.tokenClient == nil {
		return nil
	}
	return []TokenSourceOp{TokenSourceWithClient(o.tokenClient)}
}

// ClientOps creates options for functions like NewClient
func (o *RequestFlowConfig) ClientOps() []ClientOp {
	if o == nil || o.requestClient == nil {
		return nil
	}
	return []ClientOp{ClientWithRequestClient(o.requestClient)}
}

// RequestFlowOp is generic option for whole oauth2 flow
type RequestFlowOp func(o *RequestFlowConfig)

// NewRequestFlowConfig creates RequestFlowConfig from RequestFlowOp.
// nil is a valid RequestFlowConfig
func NewRequestFlowConfig(ops ...RequestFlowOp) *RequestFlowConfig {
	if len(ops) == 0 {
		return nil
	}
	options := &RequestFlowConfig{}
	for _, op := range ops {
		op(options)
	}
	return options
}

// RequestFlowWithClient is an option for both token and request clients
func RequestFlowWithClient(client *http.Client) RequestFlowOp {
	return func(o *RequestFlowConfig) {
		o.tokenClient = client
		o.requestClient = client
	}
}

// RequestFlowWithTokenClient is an option for token client
func RequestFlowWithTokenClient(client *http.Client) RequestFlowOp {
	return func(o *RequestFlowConfig) {
		o.tokenClient = client
	}
}

// RequestFlowWithRequestClient is an option for request client
func RequestFlowWithRequestClient(client *http.Client) RequestFlowOp {
	return func(o *RequestFlowConfig) {
		o.requestClient = client
	}
}

// TokenSourceConfig is a config for TokenSource.
//
// Currently, provides only optional http client.
type TokenSourceConfig struct {
	client *http.Client
}

// NewTokenSourceConfig creates TokenSourceConfig from TokenSourceOp.
// nil is a valid TokenSourceConfig
func NewTokenSourceConfig(ops ...TokenSourceOp) *TokenSourceConfig {
	if len(ops) == 0 {
		return nil
	}
	options := &TokenSourceConfig{}
	for _, op := range ops {
		op(options)
	}
	return options
}

// WithOauth2HTTPClient is used to set oauth2.HTTPClient to context
func (o *TokenSourceConfig) WithOauth2HTTPClient(ctx context.Context) context.Context {
	if o == nil || o.client == nil {
		return ctx
	}
	return context.WithValue(ctx, oauth2.HTTPClient, o.client)
}

// GetOptionalClient returns http client or nil
func (o *TokenSourceConfig) GetOptionalClient() *http.Client {
	if o == nil {
		return nil
	}
	return o.client
}

// TokenSourceOp is an option for TokenSource
type TokenSourceOp func(o *TokenSourceConfig)

// TokenSourceWithClient is an option for token client
func TokenSourceWithClient(client *http.Client) TokenSourceOp {
	return func(o *TokenSourceConfig) {
		o.client = client
	}
}

// Oauth2ContextClient is a helper function to get oauth2.HTTPClient. Returns nil on missing value
func Oauth2ContextClient(ctx context.Context) *http.Client {
	if ctx == nil {
		return nil
	}
	if hc, ok := ctx.Value(oauth2.HTTPClient).(*http.Client); ok {
		return hc
	}
	return nil
}
