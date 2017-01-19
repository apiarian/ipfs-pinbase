package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"

	"github.com/pkg/errors"
)

const (
	JWTSubjectContextKey    = "jwt-subject"
	AuthorizationContextKey = "subject-permissions"
)

type AuthenticationRequirer struct {
	key []byte
}

func NewAuthenticationRequirer(jwtKey []byte) (*AuthenticationRequirer, error) {
	return &AuthenticationRequirer{
		key: jwtKey,
	}, nil
}

func (ar *AuthenticationRequirer) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		subject, err := ar.checkAuthHeader(r)
		if err != nil {
			log.Printf("error checking auth header: %+v", err)

			SetContentTypeJSON(w)

			MarshalResponse(
				w,
				http.StatusUnauthorized,
				ErrorResponse{
					Message: "authorization required",
					Details: "valid Bearer JWT authorization header required",
				},
			)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), JWTSubjectContextKey, subject))

		h.ServeHTTP(w, r)
	})
}

func (ar *AuthenticationRequirer) checkAuthHeader(r *http.Request) (string, error) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return "", errors.New("no authorization header")
	}

	if !strings.HasPrefix(h, "Bearer ") {
		return "", errors.New("no 'Bearer ' prefix on authorization header")
	}

	raw := strings.TrimPrefix(h, "Bearer ")

	tok, err := jwt.ParseSigned(raw)
	if err != nil {
		return "", errors.Wrap(err, "parse raw header")
	}

	if len(tok.Headers) != 1 {
		return "", errors.New("malformed token: too many headers")
	}

	if tok.Headers[0].Algorithm != string(jose.HS256) {
		return "", errors.Errorf("malformed token: incorrect algorithm: %s", tok.Headers[0].Algorithm)
	}

	var claims jwt.Claims
	if err := tok.Claims(ar.key, &claims); err != nil {
		return "", errors.Wrap(err, "extract claims")
	}

	if claims.Issuer != "pinbased" {
		return "", errors.Errorf("malformed token: incorrect issuer: %s", claims.Issuer)
	}

	if claims.Subject == "" {
		return "", errors.Errorf("malformed token: missing subject")
	}

	if claims.IssuedAt.Time().After(time.Now()) {
		return "", errors.Errorf("malformed token: claims seem to be issued from the future: %s", claims.IssuedAt.Time())
	}

	return claims.Subject, nil
}

type Authorization struct {
	Admin bool
}

type AuthorizationStorage interface {
	GetAuthorization(string) (*Authorization, error)
}

type AuthorizationInjector struct {
	store AuthorizationStorage
}

func NewAuthorizationInjector(as AuthorizationStorage) (*AuthorizationInjector, error) {
	return &AuthorizationInjector{
		store: as,
	}, nil
}

func (ai *AuthorizationInjector) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization, err := ai.getAuthorization(r)
		if err != nil {
			log.Printf("error getting authorization: %+v", err)

			SetContentTypeJSON(w)
			InternalServerError(w)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), AuthorizationContextKey, authorization))

		h.ServeHTTP(w, r)
	})
}

func (ai *AuthorizationInjector) getAuthorization(r *http.Request) (*Authorization, error) {
	subject, ok := r.Context().Value(JWTSubjectContextKey).(string)

	if !ok {
		return nil, errors.New("no JWT subject in context")
	}

	auth, err := ai.store.GetAuthorization(subject)

	return auth, errors.Wrap(err, "get authorization")
}
