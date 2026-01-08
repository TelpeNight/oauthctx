package oauthctx

import (
	"context"

	"golang.org/x/oauth2"
)

// Oauth2TokenSourceWithContext is an interface, implemented by some oauth2 types like clientcredentials.Config, that can be adopted by this library
type Oauth2TokenSourceWithContext interface {
	Token(ctx context.Context) (*oauth2.Token, error)
}

// NewOauth2Token is a factory func, which implements Oauth2TokenSourceWithContext
type NewOauth2Token func(ctx context.Context) (*oauth2.Token, error)

// Token implements Oauth2TokenSourceWithContext
func (n NewOauth2Token) Token(ctx context.Context) (*oauth2.Token, error) {
	return n(ctx)
}

// AdoptTokenSourceWithContext converts Oauth2TokenSourceWithContext to TokenSource
func AdoptTokenSourceWithContext(src Oauth2TokenSourceWithContext, ops ...TokenSourceOp) TokenSource {
	return &adoptedTokenSourcer{
		new:  src,
		conf: NewTokenSourceConfig(ops...),
	}
}

// Oauth2TokenConfig is an interface, implemented by some oauth2 types (usually configs), that can be adopted by this library
type Oauth2TokenConfig interface {
	TokenSource(ctx context.Context) oauth2.TokenSource
}

// NewOauth2TokenSource is a factory func, which implements Oauth2TokenConfig
type NewOauth2TokenSource func(ctx context.Context) oauth2.TokenSource

// TokenSource implements Oauth2TokenConfig
func (n NewOauth2TokenSource) TokenSource(ctx context.Context) oauth2.TokenSource {
	return n(ctx)
}

// AdoptTokenConfig converts Oauth2TokenConfig to TokenSource
func AdoptTokenConfig(src Oauth2TokenConfig, ops ...TokenSourceOp) TokenSource {
	return &adoptedTokenSourcer{
		new: NewOauth2Token(func(ctx context.Context) (*oauth2.Token, error) {
			return src.TokenSource(ctx).Token()
		}),
		conf: NewTokenSourceConfig(ops...),
	}
}

type adoptedTokenSourcer struct {
	new  Oauth2TokenSourceWithContext
	conf *TokenSourceConfig
}

func (c *adoptedTokenSourcer) TokenContext(ctx context.Context) (*oauth2.Token, error) {
	return c.new.Token(c.conf.WithOauth2HTTPClient(ctx))
}
