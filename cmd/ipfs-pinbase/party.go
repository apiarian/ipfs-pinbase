package main

import (
	"github.com/apiarian/ipfs-pinbase/cmd/ipfs-pinbase/app"
	"github.com/goadesign/goa"
)

// PartyController implements the party resource.
type PartyController struct {
	*goa.Controller
}

// NewPartyController creates a party controller.
func NewPartyController(service *goa.Service) *PartyController {
	return &PartyController{Controller: service.NewController("PartyController")}
}

// Create runs the create action.
func (c *PartyController) Create(ctx *app.CreatePartyContext) error {
	// PartyController_Create: start_implement

	// Put your logic here

	// PartyController_Create: end_implement
	return nil
}

// Delete runs the delete action.
func (c *PartyController) Delete(ctx *app.DeletePartyContext) error {
	// PartyController_Delete: start_implement

	// Put your logic here

	// PartyController_Delete: end_implement
	return nil
}

// List runs the list action.
func (c *PartyController) List(ctx *app.ListPartyContext) error {
	// PartyController_List: start_implement

	// Put your logic here

	// PartyController_List: end_implement
	res := app.PinbasePartyCollection{}
	return ctx.OK(res)
}

// Show runs the show action.
func (c *PartyController) Show(ctx *app.ShowPartyContext) error {
	// PartyController_Show: start_implement

	// Put your logic here

	// PartyController_Show: end_implement
	res := &app.PinbaseParty{}
	return ctx.OK(res)
}

// Update runs the update action.
func (c *PartyController) Update(ctx *app.UpdatePartyContext) error {
	// PartyController_Update: start_implement

	// Put your logic here

	// PartyController_Update: end_implement
	res := &app.PinbaseParty{}
	return ctx.OK(res)
}
