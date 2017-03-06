package main

import (
	"github.com/apiarian/ipfs-pinbase/app"
	"github.com/goadesign/goa"
	"github.com/pkg/errors"
)

// NodeController implements the node resource.
type NodeController struct {
	*goa.Controller
}

// NewNodeController creates a node controller.
func NewNodeController(service *goa.Service) *NodeController {
	return &NodeController{Controller: service.NewController("NodeController")}
}

// Create runs the create action.
func (c *NodeController) Create(ctx *app.CreateNodeContext) error {
	// NodeController_Create: start_implement

	// Put your logic here

	// NodeController_Create: end_implement
	return nil
}

// Delete runs the delete action.
func (c *NodeController) Delete(ctx *app.DeleteNodeContext) error {
	// NodeController_Delete: start_implement

	// Put your logic here

	// NodeController_Delete: end_implement
	return nil
}

// List runs the list action.
func (c *NodeController) List(ctx *app.ListNodeContext) error {
	// NodeController_List: start_implement

	subject, err := ExtractJWTSubject(ctx)
	if err != nil {
		return errors.Wrap(err, "extract JWT subject")
	}

	goa.LogInfo(ctx, "found a JWT", "subject", subject)

	// NodeController_List: end_implement
	res := app.PinbaseNodeCollection{}
	return ctx.OK(res)
}

// Show runs the show action.
func (c *NodeController) Show(ctx *app.ShowNodeContext) error {
	// NodeController_Show: start_implement

	// Put your logic here

	// NodeController_Show: end_implement
	res := &app.PinbaseNode{}
	return ctx.OK(res)
}

// Update runs the update action.
func (c *NodeController) Update(ctx *app.UpdateNodeContext) error {
	// NodeController_Update: start_implement

	// Put your logic here

	// NodeController_Update: end_implement
	res := &app.PinbaseNode{}
	return ctx.OK(res)
}
