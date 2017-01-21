package main

import (
	"github.com/apiarian/ipfs-pinbase/app"
	"github.com/goadesign/goa"
)

// NodeController implements the node resource.
type NodeController struct {
	*goa.Controller
}

// NewNodeController creates a node controller.
func NewNodeController(service *goa.Service) *NodeController {
	return &NodeController{Controller: service.NewController("NodeController")}
}

// Show runs the show action.
func (c *NodeController) Show(ctx *app.ShowNodeContext) error {
	// NodeController_Show: start_implement

	// Put your logic here

	// NodeController_Show: end_implement
	res := &app.PinbaseNode{
		Hash:        ctx.NodeHash,
		Description: "something about " + ctx.NodeHash,
		APIURL:      "https://foobar.foobar:12345",
	}
	return ctx.OK(res)
}
