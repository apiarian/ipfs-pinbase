package main

import (
	"github.com/apiarian/ipfs-pinbase/_scaffolds/app"
	"github.com/goadesign/goa"
)

// LoginController implements the login resource.
type LoginController struct {
	*goa.Controller
}

// NewLoginController creates a login controller.
func NewLoginController(service *goa.Service) *LoginController {
	return &LoginController{Controller: service.NewController("LoginController")}
}

// Login runs the login action.
func (c *LoginController) Login(ctx *app.LoginLoginContext) error {
	// LoginController_Login: start_implement

	// Put your logic here

	// LoginController_Login: end_implement
	return nil
}
