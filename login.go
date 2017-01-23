package main

import (
	"time"

	"github.com/apiarian/ipfs-pinbase/app"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa"
	"github.com/pkg/errors"
)

// LoginController implements the login resource.
type LoginController struct {
	*goa.Controller
	jwtKey []byte
}

// NewLoginController creates a login controller.
func NewLoginController(service *goa.Service, jwtKey []byte) *LoginController {
	return &LoginController{
		Controller: service.NewController("LoginController"),
		jwtKey:     jwtKey,
	}
}

// Login runs the login action.
func (c *LoginController) Login(ctx *app.LoginLoginContext) error {
	type PinbaseClaims struct {
		Scopes []string `json:"scopes"`
		jwtgo.StandardClaims
	}

	cl := PinbaseClaims{
		Scopes: []string{
			"node:view",
			"node:edit",
		},
		StandardClaims: jwtgo.StandardClaims{
			Issuer:   "pinbased",
			Subject:  "someuser",
			IssuedAt: time.Now().Unix(),
		},
	}
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, cl)

	signed, err := token.SignedString(c.jwtKey)
	if err != nil {
		return errors.Wrap(err, "sign token")
	}

	ctx.ResponseData.Header().Set(
		"Authorization",
		"Bearer "+signed,
	)

	return ctx.NoContent()
}
