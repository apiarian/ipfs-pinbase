package main

import (
	"github.com/apiarian/ipfs-pinbase/cmd/ipfs-pinbase/_scaffolds/app"
	"github.com/goadesign/goa"
)

// PinController implements the pin resource.
type PinController struct {
	*goa.Controller
}

// NewPinController creates a pin controller.
func NewPinController(service *goa.Service) *PinController {
	return &PinController{Controller: service.NewController("PinController")}
}

// Create runs the create action.
func (c *PinController) Create(ctx *app.CreatePinContext) error {
	// PinController_Create: start_implement

	// Put your logic here

	// PinController_Create: end_implement
	return nil
}

// Delete runs the delete action.
func (c *PinController) Delete(ctx *app.DeletePinContext) error {
	// PinController_Delete: start_implement

	// Put your logic here

	// PinController_Delete: end_implement
	return nil
}

// List runs the list action.
func (c *PinController) List(ctx *app.ListPinContext) error {
	// PinController_List: start_implement

	// Put your logic here

	// PinController_List: end_implement
	res := app.PinbasePinCollection{}
	return ctx.OK(res)
}

// Show runs the show action.
func (c *PinController) Show(ctx *app.ShowPinContext) error {
	// PinController_Show: start_implement

	// Put your logic here

	// PinController_Show: end_implement
	res := &app.PinbasePin{}
	return ctx.OK(res)
}

// Update runs the update action.
func (c *PinController) Update(ctx *app.UpdatePinContext) error {
	// PinController_Update: start_implement

	// Put your logic here

	// PinController_Update: end_implement
	res := &app.PinbasePin{}
	return ctx.OK(res)
}
