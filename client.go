package oauthctx

import (
	"net/http"
)

// NewClient creates an *http.Client from TokenSource.
func NewClient(src TokenSource, ops ...ClientOp) *http.Client {
	var options clientOp
	for _, op := range ops {
		op(&options)
	}
	if options.client == nil {
		options.client = http.DefaultClient
	}
	if src == nil {
		// to match original behaviour
		return options.client
	}
	return &http.Client{
		Transport: &Transport{
			Source: ReuseTokenSource(nil, src),
			Base:   options.client.Transport,
		},
		CheckRedirect: options.client.CheckRedirect,
		Jar:           options.client.Jar,
		Timeout:       options.client.Timeout,
	}
}

// ClientOp is an option for NewClient
type ClientOp func(o *clientOp)

// ClientWithRequestClient is an option for underlying request client
func ClientWithRequestClient(client *http.Client) ClientOp {
	return func(o *clientOp) {
		o.client = client
	}
}

type clientOp struct {
	client *http.Client
}
