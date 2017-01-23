package main

import (
	"net/http"

	"golang.org/x/net/context"

	"github.com/apiarian/ipfs-pinbase/app"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware/security/basicauth"
	"github.com/goadesign/goa/middleware/security/jwt"
	"github.com/pkg/errors"
)

func NewBasicAuth() (goa.Middleware, error) {
	return func(h goa.Handler) goa.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			u, p, ok := r.BasicAuth()
			if !ok || u != "user" || p != "kittens" {
				return basicauth.ErrBasicAuthFailed("Authentication failed")
			}

			return h(ctx, w, r)
		}
	}, nil
}

func NewJWTAuth(jwtKey []byte) (goa.Middleware, error) {
	return jwt.New(
		jwt.NewSimpleResolver(
			[]jwt.Key{jwtKey},
		),
		func(h goa.Handler) goa.Handler {
			return h
			/*
				return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
					return h(ctx, w, r)
				}
			*/
		},
		app.NewJWTSecurity(),
	), nil
}

func ExtractJWTSubject(ctx context.Context) (string, error) {
	token := jwt.ContextJWT(ctx)
	if token == nil {
		return "", errors.New("no JWT in context")
	}

	claims := token.Claims.(jwtgo.MapClaims)

	return claims["sub"].(string), nil
}
