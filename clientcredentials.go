package oauthctx

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// ClientCredentials describes a 2-legged OAuth2 flow, with both the
// client application information and the server's endpoint URLs.
type ClientCredentials struct{ *clientcredentials.Config }

// NewClientCredentials creates new ClientCredentials
func NewClientCredentials(cfg *clientcredentials.Config) *ClientCredentials {
	return &ClientCredentials{cfg}
}

// Client returns an HTTP client using the provided token.
// The token will auto-refresh as necessary.
//
// The provided options control which HTTP client
// is used.
//
// The returned Client and its Transport should not be modified.
func (c *ClientCredentials) Client(ops ...RequestFlowOp) *http.Client {
	conf := NewRequestFlowConfig(ops...)
	return NewClient(
		c.tokenSource(conf.TokenSourceOps()), // NewClient will reuse tokenSource
		conf.ClientOps()...)
}

// Token uses client credentials to retrieve a token.
func (c *ClientCredentials) Token(ctx context.Context) (*oauth2.Token, error) {
	return c.Config.Token(ctx)
}

// TokenWithOptions uses client credentials to retrieve a token.
//
// The provided options optionally controls which HTTP client is used.
func (c *ClientCredentials) TokenWithOptions(ctx context.Context, ops ...TokenSourceOp) (*oauth2.Token, error) {
	return c.tokenSource(ops).TokenContext(ctx)
}

// TokenSource returns a TokenSource that returns t until t expires,
// automatically refreshing it as necessary using the provided options and the
// client ID and client secret.
//
// Most users will use Config.Client instead.
func (c *ClientCredentials) TokenSource(ops ...TokenSourceOp) TokenSource {
	return ReuseTokenSource(nil, c.tokenSource(ops))
}

func (c *ClientCredentials) tokenSource(ops []TokenSourceOp) TokenSource {
	return AdoptTokenSourceWithContext(c.Config, ops...)
}
