package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type fakeAuther struct{}

func (f fakeAuther) CheckPassword(username, password string) (bool, error) {
	switch username {
	case "realuser":
		return (password == "correctpassword"), nil

	case "backendfailure":
		return false, fmt.Errorf("authentication back-end failure")

	default:
		panic("only realuser or backendfailure usernames are allowed")
	}
}

func TestLogin(t *testing.T) {
	key := []byte("secretkey")

	login, err := NewLogin(
		key,
		fakeAuther{},
	)
	if err != nil {
		t.Fatalf("failed to create login: %+v", err)
	}

	cases := []struct {
		tag     string
		prep    func(*http.Request)
		code    int
		content map[string]interface{}
	}{
		{
			tag: "correct login",
			prep: func(r *http.Request) {
				r.SetBasicAuth("realuser", "correctpassword")
			},
			code:    http.StatusOK,
			content: nil,
		},
		{
			tag:  "missing header",
			prep: func(r *http.Request) {},
			code: http.StatusBadRequest,
			content: map[string]interface{}{
				"error-message": "missing auth",
				"error-details": "missing Authorization header",
			},
		},
		{
			tag: "bad password",
			prep: func(r *http.Request) {
				r.SetBasicAuth("realuser", "incorrectpassword")
			},
			code: http.StatusUnauthorized,
			content: map[string]interface{}{
				"error-message": "authentication required",
				"error-details": "login failed",
			},
		},
		{
			tag: "backend failure",
			prep: func(r *http.Request) {
				r.SetBasicAuth("backendfailure", "doesnotmatter")
			},
			code: http.StatusInternalServerError,
			content: map[string]interface{}{
				"error-message": "internal server error",
				"error-details": "can't even.",
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.tag, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodPost, "/login", nil)
			c.prep(r)

			w := httptest.NewRecorder()

			login.ServeHTTP(w, r)

			if want, got := "application/json", w.Header().Get("Content-Type"); got != want {
				t.Errorf("content-type: wanted %s, got %s", want, got)
			}

			if want, got := c.code, w.Code; got != want {
				t.Errorf("error code: wanted %d, got %d", want, got)
			}

			var d map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&d)
			if err != nil {
				t.Errorf("failed to decode response JSON: %+v", err)
			}

			if c.content != nil {
				if !reflect.DeepEqual(c.content, d) {
					t.Errorf("response content: wanted %+v, got %+v", c.content, d)
				}
			} else {
				if want, got := 1, len(d); got != want {
					t.Errorf("number of keys in data: want %d, got %d (%+v)", want, got, d)
				}

				raw, ok := d["jwt"].(string)
				if !ok {
					t.Errorf("no jwt string field in data: %+v", d)
				}

				tok, err := jwt.ParseSigned(raw)
				if err != nil {
					t.Errorf("failed to parse token: %+v", err)
				}

				if want, got := 1, len(tok.Headers); got != want {
					t.Errorf("number of token headers: want %d, got %d", want, got)
				}

				if want, got := jose.HS256, tok.Headers[0].Algorithm; string(want) != got {
					t.Errorf("JWT algorithm: want %s, got %s", want, got)
				}

				var claims jwt.Claims
				if err := tok.Claims(key, &claims); err != nil {
					t.Errorf("failed to extract claims: %+v", err)
				}

				if want, got := "pinbased", claims.Issuer; got != want {
					t.Errorf("iss (Issuer): want %s, got %s", want, got)
				}

				if want, got := "realuser", claims.Subject; got != want {
					t.Errorf("sub (Subject): want %s, got %s", want, got)
				}

				n := time.Now()
				if claims.IssuedAt.Time().Sub(n) > 1*time.Microsecond {
					t.Errorf("iss (IssuedAt) more than a second in the past: %s (now %s)", claims.IssuedAt.Time(), n)
				}
			}
		})
	}
}
