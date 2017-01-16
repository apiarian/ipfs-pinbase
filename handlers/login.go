package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type Authenticator interface {
	CheckPassword(username, password string) (bool, error)
}

type Login struct {
	signer jose.Signer
	auther Authenticator
}

func NewLogin(jwtKey []byte, auther Authenticator) (*Login, error) {
	sig, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.HS256,
			Key:       jwtKey,
		},
		&jose.SignerOptions{},
	)
	if err != nil {
		return nil, errors.Wrap(err, "create signer for JWT")
	}

	return &Login{
		signer: sig,
		auther: auther,
	}, nil
}

type LoginResponse struct {
	JWT string `json:"jwt"`
}

type ErrorResponse struct {
	Message string `json:"error-message"`
	Details string `json:"error-details"`
}

func MarshalResponse(w http.ResponseWriter, c int, data interface{}) {
	j, err := json.Marshal(data)

	if err != nil {
		log.Printf("failed to marshal JSON: %+v", err)

		w.WriteHeader(http.StatusInternalServerError)

		w.Write(
			[]byte(
				`{"error-message":"internal server error","error-details":"can't even."}`,
			),
		)

		return
	}

	w.WriteHeader(c)
	_, err = w.Write(j)

	if err != nil {
		log.Printf("failed to write response: %+v", err)
	}
}

func InternalServerError(w http.ResponseWriter) {
	MarshalResponse(
		w,
		http.StatusInternalServerError,
		ErrorResponse{
			Message: "internal server error",
			Details: "can't even.",
		},
	)
}

func (l *Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	username, password, ok := r.BasicAuth()

	if !ok {
		MarshalResponse(
			w,
			http.StatusBadRequest,
			ErrorResponse{
				Message: "missing auth",
				Details: "missing Authorization header",
			},
		)
		return
	}

	authed, err := l.auther.CheckPassword(username, password)
	if err != nil {
		log.Printf("failed to check password: %+v", err)

		InternalServerError(w)
		return
	}
	if !authed {
		MarshalResponse(
			w,
			http.StatusUnauthorized,
			ErrorResponse{
				Message: "authentication required",
				Details: "login failed",
			},
		)
		return
	}

	cl := jwt.Claims{
		Issuer:   "pinbased",
		Subject:  username,
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}

	raw, err := jwt.Signed(l.signer).Claims(cl).CompactSerialize()
	if err != nil {
		log.Printf("failed to create signed token: %+v", err)

		InternalServerError(w)
		return
	}

	MarshalResponse(
		w,
		http.StatusOK,
		LoginResponse{
			JWT: raw,
		},
	)
}
