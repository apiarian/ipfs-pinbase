package main

import (
	"time"

	"github.com/apiarian/ipfs-pinbase/app"
	"github.com/goadesign/goa"
	"github.com/pkg/errors"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// LoginController implements the login resource.
type LoginController struct {
	*goa.Controller
	signer jose.Signer
}

// NewLoginController creates a login controller.
func NewLoginController(service *goa.Service, jwtKey []byte) *LoginController {
	sig, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.HS256,
			Key:       jwtKey,
		},
		&jose.SignerOptions{},
	)
	if err != nil {
		panic(err)
	}

	return &LoginController{
		Controller: service.NewController("LoginController"),
		signer:     sig,
	}
}

type ExtraClaims struct {
	Scopes []string `json:"scopes"`
}

// Login runs the login action.
func (c *LoginController) Login(ctx *app.LoginLoginContext) error {
	cl := jwt.Claims{
		Issuer:   "pinbased",
		Subject:  "someuser",
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}

	extra := ExtraClaims{
		Scopes: []string{
			"node:view",
			"node:edit",
		},
	}

	raw, err := jwt.Signed(c.signer).Claims(cl).Claims(extra).CompactSerialize()
	if err != nil {
		return errors.Wrap(err, "sign token")
	}

	ctx.ResponseData.Header().Set(
		"Authorization",
		"Bearer "+raw,
	)

	return ctx.NoContent()
}
