[![Go Reference](https://pkg.go.dev/badge/github.com/TelpeNight/oauthctx.svg)](https://pkg.go.dev/github.com/TelpeNight/oauthctx)

# oauthctx

---

This library is used to bypass golang oauth2 package token source limitations: https://github.com/golang/oauth2/issues/262

The aim ot this library is to implement only the smallest possible subset of functionality for context support.
It will reuse original library as match as possible. This goal is achieved through these steps:

1. Reuse basic token retrieval by adopting existing data types, which implement `Oauth2TokenConfig` and `Oauth2TokenSourceWithContext`.
See `convert.go`
2. Reusing high-level functions by retrieving Token through context-aware `TokenSource` and passing it to the existing implementation with the help of `oauth2.StaticTokenSourc`.
See `transport.go` and `grpc/credentials.go`
3. Reimplement only a small subset like `ReuseTokenSource` or `tokenRefresher`, which is focused and bug-free.

All configuration is provided with the option pattern (no more context.WithValue mess).

All functions, provided in this library, should be familiar to oauth2 users.
If something is missing - feel free to open an issue or make a PR.
I believe, that everything can be adopted by this approach.
The library is already used in production.

## Code examples

---

```go
// grpc
package main
import (
    "golang.org/x/oauth2"
	
    "google.golang.org/grpc/credentials"
    gcred "google.golang.org/grpc/credentials/google"

    "github.com/TelpeNight/oauthctx"
    grpcctx "github.com/TelpeNight/oauthctx/grpc"
)

var conf = oauthctx.NewConfig(&oauth2.Config{
    //...
})
var refreshToken string = "..."

ts := conf.TokenSource(
	&oauth2.Token{RefreshToken: refreshToken},
	// custom http.Client can be provided with option
	oauthctx.TokenSourceWithClient(...))
ts = oauthctx.ReuseTokenSource(nil, ts)

var bundle credentials.Bundle = gcred.NewDefaultCredentialsWithOptions(
    gcred.DefaultCredentialsOptions{
        PerRPCCreds: &grpcctx.TokenSource{
            TokenSource: ts,
        },
    },
)

// use bundle to create a client. methods' context will be passed to oauth2, so overall call will respect timeouts
```
