package main

import (
	"github.com/apiarian/ipfs-pinbase/cmd/ipfs-pinbase/app"
	"github.com/apiarian/ipfs-pinbase/pinbase"
	"github.com/goadesign/goa"
)

// PartyController implements the party resource.
type PartyController struct {
	*goa.Controller
	P pinbase.PinProvider
}

// NewPartyController creates a party controller.
func NewPartyController(service *goa.Service, P pinbase.PinProvider) *PartyController {
	return &PartyController{Controller: service.NewController("PartyController"), P: P}
}

// Create runs the create action.
func (c *PartyController) Create(ctx *app.CreatePartyContext) error {
	// PartyController_Create: start_implement

	err := c.P.PinService().CreateParty(&pinbase.PartyCreate{
		ID:          pinbase.Hash(ctx.Payload.Hash),
		Description: ctx.Payload.Description,
	})
	if err != nil {
		return err
	}

	// PartyController_Create: end_implement
	return nil
}

// Delete runs the delete action.
func (c *PartyController) Delete(ctx *app.DeletePartyContext) error {
	// PartyController_Delete: start_implement

	err := c.P.PinService().DeleteParty(pinbase.Hash(ctx.PartyHash))
	if err != nil {
		return err
	}

	// PartyController_Delete: end_implement
	return nil
}

// List runs the list action.
func (c *PartyController) List(ctx *app.ListPartyContext) error {
	// PartyController_List: start_implement

	p, err := c.P.PinService().Parties()
	if err != nil {
		return err
	}

	res := app.PinbasePartyCollection{}
	for _, x := range p {
		res = append(res, &app.PinbaseParty{
			Hash:        string(x.ID),
			Description: x.Description,
		})
	}

	// PartyController_List: end_implement
	return ctx.OK(res)
}

// Show runs the show action.
func (c *PartyController) Show(ctx *app.ShowPartyContext) error {
	// PartyController_Show: start_implement

	p, err := c.P.PinService().Party(pinbase.Hash(ctx.PartyHash))
	if err != nil {
		return err
	}

	res := &app.PinbaseParty{
		Hash:        string(p.ID),
		Description: p.Description,
	}

	// PartyController_Show: end_implement
	return ctx.OK(res)
}

// Update runs the update action.
func (c *PartyController) Update(ctx *app.UpdatePartyContext) error {
	// PartyController_Update: start_implement

	ps := c.P.PinService()

	err := ps.UpdateParty(
		pinbase.Hash(ctx.PartyHash),
		&pinbase.PartyEdit{
			Description: *ctx.Payload.Description,
		},
	)
	if err != nil {
		return err
	}

	p, err := ps.Party(pinbase.Hash(ctx.PartyHash))
	if err != nil {
		return err
	}

	res := &app.PinbaseParty{
		Hash:        string(p.ID),
		Description: p.Description,
	}

	// PartyController_Update: end_implement
	return ctx.OK(res)
}
