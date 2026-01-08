package oauthctx

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

// Config describes a typical 3-legged OAuth2 flow, with both the
// client application information and the server's endpoint URLs.
// For the client credentials 2-legged OAuth2 flow, see the ClientCredentials.
type Config struct{ *oauth2.Config }

// NewConfig creates new Config
func NewConfig(cfg *oauth2.Config) *Config {
	return &Config{Config: cfg}
}

// Client returns an HTTP client using the provided token.
// The token will auto-refresh as necessary. The underlying
// HTTP transport will be obtained from options.
// The returned client and its Transport should not be modified.
func (c *Config) Client(t *oauth2.Token, ops ...RequestFlowOp) *http.Client {
	conf := NewRequestFlowConfig(ops...)
	return NewClient(
		c.tokenSource(t, conf.TokenSourceOps()), // NewClient will reuse tokenSource
		conf.ClientOps()...)
}

// TokenSource returns a TokenSource that returns t until t expires,
// automatically refreshing it as necessary.
//
// Most users will use Config.Client instead.
func (c *Config) TokenSource(t *oauth2.Token, ops ...TokenSourceOp) TokenSource {
	src := c.tokenSource(t, ops)
	return ReuseTokenSource(t, src)
}

func (c *Config) tokenSource(t *oauth2.Token, ops []TokenSourceOp) TokenSource {
	tkr := &tokenRefresher{
		srcConf:   NewTokenSourceConfig(ops...),
		oauthConf: c.Config,
	}
	if t != nil {
		tkr.refreshToken = t.RefreshToken
	}
	return tkr
}

// tokenRefresher is a TokenSource that makes "grant_type"=="refresh_token"
// HTTP requests to renew a token using a RefreshToken.
type tokenRefresher struct {
	srcConf      *TokenSourceConfig
	oauthConf    *oauth2.Config
	refreshToken string
}

// TokenContext WARNING: is not safe for concurrent access, as it
// updates the tokenRefresher's refreshToken field.
// Within this package, it is used by reuseTokenSource which
// synchronizes calls to this method with its own mutex.
func (tf *tokenRefresher) TokenContext(ctx context.Context) (*oauth2.Token, error) {
	src := tf.oauthConf.TokenSource(
		tf.srcConf.WithOauth2HTTPClient(ctx),
		&oauth2.Token{RefreshToken: tf.refreshToken},
	)
	tk, err := src.Token()
	if err != nil {
		return nil, err
	}
	tf.refreshToken = tk.RefreshToken
	return tk, nil
}
