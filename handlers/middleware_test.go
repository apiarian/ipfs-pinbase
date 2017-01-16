package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

func makeSignedPrep(sig jose.Signer, c jwt.Claims) func(*http.Request) {
	return func(r *http.Request) {
		raw, err := jwt.Signed(sig).Claims(c).CompactSerialize()
		if err != nil {
			panic(err)
		}

		r.Header.Set("Authorization", "Bearer "+raw)
	}
}

func TestCheckAuthHeader(t *testing.T) {
	key := []byte("secretkey")

	sig, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.HS256,
			Key:       key,
		},
		&jose.SignerOptions{},
	)
	if err != nil {
		t.Fatalf("failed to create JOSE signer: %+v", err)
	}

	wrongSig, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.HS512,
			Key:       key,
		},
		&jose.SignerOptions{},
	)
	if err != nil {
		t.Fatalf("faield to create wrong signing JOSE signer: %+v", err)
	}

	altSig, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.HS256,
			Key:       []byte("otherkey"),
		},
		&jose.SignerOptions{},
	)

	auther, err := NewAuthRequirer(key)
	if err != nil {
		t.Fatalf("failed to create auth requirer: %+v", err)
	}

	cases := []struct {
		tag     string
		prep    func(*http.Request)
		err     error
		subject string
	}{
		{
			tag: "good header",
			prep: makeSignedPrep(
				sig,
				jwt.Claims{
					Issuer:   "pinbased",
					Subject:  "realuser",
					IssuedAt: jwt.NewNumericDate(time.Now()),
				},
			),
			err:     nil,
			subject: "realuser",
		},
		{
			tag:     "missing header",
			prep:    func(r *http.Request) {},
			err:     errors.New("no authorization header"),
			subject: "",
		},
		{
			tag: "missing bearer prefix",
			prep: func(r *http.Request) {
				r.Header.Set("Authorization", "foobar")
			},
			err:     errors.New("no 'Bearer ' prefix on authorization header"),
			subject: "",
		},
		{
			tag: "malformed jwt",
			prep: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer foobar")
			},
			err:     errors.New("parse raw header: "),
			subject: "",
		},
		/*
			// TODO: find a way to put in multiple headers
			{
				tag: "multiple headers",
				prep: func(r *http.Request) {
				},
				err:     errors.New("malformed token: too many headers"),
				subject: "",
			},
		*/
		{
			tag: "wrong algorithm",
			prep: makeSignedPrep(
				wrongSig,
				jwt.Claims{
					Issuer:   "pinbased",
					Subject:  "realuser",
					IssuedAt: jwt.NewNumericDate(time.Now()),
				},
			),
			err:     errors.New("malformed token: incorrect algorithm: HS512"),
			subject: "",
		},
		{
			tag: "wrong key",
			prep: makeSignedPrep(
				altSig,
				jwt.Claims{
					Issuer:   "pinbased",
					Subject:  "realuser",
					IssuedAt: jwt.NewNumericDate(time.Now()),
				},
			),
			err:     errors.New("extract claims: "),
			subject: "",
		},
		{
			tag: "wring issuer",
			prep: makeSignedPrep(
				sig,
				jwt.Claims{
					Issuer:   "joebob",
					Subject:  "realuser",
					IssuedAt: jwt.NewNumericDate(time.Now()),
				},
			),
			err:     errors.New("malformed token: incorrect issuer: joebob"),
			subject: "",
		},
		{
			tag: "missing subject",
			prep: makeSignedPrep(
				sig,
				jwt.Claims{
					Issuer:   "pinbased",
					IssuedAt: jwt.NewNumericDate(time.Now()),
				},
			),
			err:     errors.New("malformed token: missing subject"),
			subject: "",
		},
		{
			tag: "jwt from the future",
			prep: makeSignedPrep(
				sig,
				jwt.Claims{
					Issuer:   "pinbased",
					Subject:  "realuser",
					IssuedAt: jwt.NewNumericDate(time.Now().Add(10 * time.Second)),
				},
			),
			err:     errors.New("malformed token: claims seem to be issued from the future: "),
			subject: "",
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.tag, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodGet, "/foo", nil)
			c.prep(r)

			subject, err := auther.checkAuthHeader(r)

			if err == nil {
				if c.err != nil {
					t.Error("did not get the expected error")
				}
			} else {
				if c.err == nil {
					t.Errorf("got an unexpected error: %+v", err)
				} else {
					if want, got := c.err, err; !strings.HasPrefix(got.Error(), want.Error()) {
						t.Errorf("error: \nwant %s...\n got %s", want, got)
					}
				}
			}

			if want, got := c.subject, subject; got != want {
				t.Errorf("subject: want %s, got %s", want, got)
			}
		})
	}
}

func TestAuthRequirerMiddlware(t *testing.T) {
	key := []byte("secretkey")

	sig, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.HS256,
			Key:       key,
		},
		&jose.SignerOptions{},
	)
	if err != nil {
		t.Fatalf("failed to create JOSE signer: %+v", err)
	}

	auther, err := NewAuthRequirer(key)
	if err != nil {
		t.Fatalf("failed to create auth requirer: %+v", err)
	}

	cases := []struct {
		tag     string
		prep    func(*http.Request)
		code    int
		content map[string]interface{}
	}{
		{
			tag: "all good",
			prep: makeSignedPrep(
				sig,
				jwt.Claims{
					Issuer:   "pinbased",
					Subject:  "realuser",
					IssuedAt: jwt.NewNumericDate(time.Now()),
				},
			),
			code: http.StatusOK,
			content: map[string]interface{}{
				"jwt-subject": "realuser",
			},
		},
		{
			tag:  "missing auth",
			prep: func(r *http.Request) {},
			code: http.StatusUnauthorized,
			content: map[string]interface{}{
				"error-message": "authorization required",
				"error-details": "valid Bearer JWT authorization header required",
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.tag, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodGet, "/foo", nil)
			c.prep(r)

			w := httptest.NewRecorder()

			wrappedHandler := auther.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				type SubjectEcho struct {
					JWTSubject string `json:"jwt-subject"`
				}

				SetContentTypeJSON(w)

				MarshalResponse(
					w,
					http.StatusOK,
					SubjectEcho{
						JWTSubject: r.Context().Value(JWTSubjectContextKey).(string),
					},
				)
				return
			}))

			wrappedHandler.ServeHTTP(w, r)

			if want, got := "application/json", w.Header().Get("Content-Type"); got != want {
				t.Errorf("content-type: wanted %s, got %s", want, got)
			}

			if want, got := c.code, w.Code; got != want {
				t.Errorf("http code: wanted %d, got %d", want, got)
			}

			var d map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&d)
			if err != nil {
				t.Errorf("failed to decode response JSON: %+v", err)
			}

			if !reflect.DeepEqual(c.content, d) {
				t.Errorf("response content: watned %+v, got %+v", c.content, d)
			}
		})
	}
}
