# gRPC Authorization

This package provides the ability to add authorization to a gRPC server. It is made to be coupled with github.com/grpc-ecosystem/go-grpc-middleware/auth.

## Background

Under the hood it is using [https://github.com/cristalhq/jwt](https://github.com/cristalhq/jwt) for all JWT functionality. The main reason this library was chosen was because [jwt.io](https://jwt.io/) reports it as supporting all [algorithms and validations](./images/cristalhq-jwt.io.png).

Currently, in this repo, the following algorithms are supported.

- HS256, HS384, HS512
- RS256, RS384, RS512

## Usage

This example will validate that the JWT was signed by the given key and not expired.

```go
package main

import (
    "net"
    "os"

    jwt "github.com/cristalhq/jwt/v3"
    grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
    grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
    "github.com/packethost/pkg/grpc/authz"
    "google.golang.org/grpc"
)

var (
    hsKey = []byte("supersecret")
)

func main() {
    // create a Config
    // at a minimum an algorithm, scope mapping (only the methods defined here will protected), and key are needed
    config := authz.NewConfig(
        jwt.HS256,
        map[string][]string{
            "/github.com.tinkerbell.pbnj.api.v1.Machine/Power": {},
        },
        authz.WithHSKey(hsKey),
        
    )

    // the AuthFunc method can then be used with as middleware with a gRPC server
    grpcServer := grpc.NewServer(
        grpc_middleware.WithUnaryServerChain(
            grpc_auth.UnaryServerInterceptor(config.AuthFunc),
        ),
    )

    listen, err := net.Listen("tcp", ":50051")
    if err != nil {
        panic(err)
    }

    if err := grpcServer.Serve(listen); err != nil {
        os.Exit(1)
    }
}
```

This example will validate that the JWT was signed by the given key, not expired, and custom scopes match the called method.

```go
package main

import (
    "encoding/json"
    "net"
    "os"

    jwt "github.com/cristalhq/jwt/v3"
    grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
    grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
    "github.com/packethost/pkg/grpc/authz"
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

var (
    hsKey = []byte("supersecret")
)

func main() {
    // create a func for validating Scopes
    scopeFunc := func(tokenClaims []byte, scopes []string) error {
        type CustomClaims struct {
            jwt.StandardClaims
            Scopes []string `json:"scopes"`
        }
        var newClaims CustomClaims
        err := json.Unmarshal(tokenClaims, &newClaims)
        if err != nil {
            return status.Errorf(codes.Unauthenticated, "access token is invalid: %s", err.Error())
        }

        if !contains(newClaims.Scopes, "write") {
            return status.Errorf(codes.PermissionDenied, "no permission to access this RPC: no matching scope found")
        }

        return nil
    }
    // create a Config
    // at a minimum an algorithm, scope mapping (only the methods defined here will protected), and a key are needed. we set the scope validation
    // and audience on this one.
    config := authz.NewConfig(
        jwt.HS256,
        map[string][]string{
            "/github.com.tinkerbell.pbnj.api.v1.Machine/Power": {"write"},
        },
        authz.WithHSKey(hsKey),
        authz.WithValidateScopeFunc(scopeFunc),
        authz.WithAudience("admin"),
    )

    // the AuthFunc method can then be used with as middleware with a gRPC server
    grpcServer := grpc.NewServer(
        grpc_middleware.WithUnaryServerChain(
            grpc_auth.UnaryServerInterceptor(config.AuthFunc),
        ),
    )

    listen, err := net.Listen("tcp", ":50051")
    if err != nil {
        panic(err)
    }

    if err := grpcServer.Serve(listen); err != nil {
        os.Exit(1)
    }
}

func contains(s []string, str string) bool {
    for _, v := range s {
        if v == str {
            return true
        }
    }
    return false
}
```
