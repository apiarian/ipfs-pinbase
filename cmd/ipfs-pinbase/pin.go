package main

import (
	"github.com/apiarian/ipfs-pinbase/cmd/ipfs-pinbase/app"
	"github.com/apiarian/ipfs-pinbase/pinbase"
	"github.com/goadesign/goa"
)

// PinController implements the pin resource.
type PinController struct {
	*goa.Controller
	P pinbase.PinProvider
}

// NewPinController creates a pin controller.
func NewPinController(service *goa.Service, P pinbase.PinProvider) *PinController {
	return &PinController{Controller: service.NewController("PinController"), P: P}
}

// Create runs the create action.
func (c *PinController) Create(ctx *app.CreatePinContext) error {
	// PinController_Create: start_implement

	err := c.P.PinService().CreatePin(
		pinbase.Hash(ctx.PartyHash),
		&pinbase.PinCreate{
			ID:         pinbase.Hash(ctx.Payload.Hash),
			Aliases:    ctx.Payload.Aliases,
			WantPinned: ctx.Payload.WantPinned,
		},
	)
	if err != nil {
		return err
	}

	// PinController_Create: end_implement
	ctx.ResponseData.Header().Set("Location", app.PinHref(ctx.PartyHash, ctx.Payload.Hash))
	return ctx.Created()
}

// Delete runs the delete action.
func (c *PinController) Delete(ctx *app.DeletePinContext) error {
	// PinController_Delete: start_implement

	err := c.P.PinService().DeletePin(
		pinbase.Hash(ctx.PartyHash),
		pinbase.Hash(ctx.PinHash),
	)
	if err != nil {
		return err
	}

	// PinController_Delete: end_implement
	return nil
}

// List runs the list action.
func (c *PinController) List(ctx *app.ListPinContext) error {
	// PinController_List: start_implement

	ps, err := c.P.PinService().Pins(pinbase.Hash(ctx.PartyHash))
	if err != nil {
		return err
	}

	res := app.PinbasePinCollection{}
	for _, p := range ps {
		var e string
		if p.LastError != nil {
			e = p.LastError.Error()
		}

		res = append(res, &app.PinbasePin{
			Hash:       string(p.ID),
			Aliases:    p.Aliases,
			WantPinned: p.WantPinned,
			Status:     p.Status.String(),
			LastError:  e,
		})
	}

	// PinController_List: end_implement
	return ctx.OK(res)
}

// Show runs the show action.
func (c *PinController) Show(ctx *app.ShowPinContext) error {
	// PinController_Show: start_implement

	p, err := c.P.PinService().Pin(
		pinbase.Hash(ctx.PartyHash),
		pinbase.Hash(ctx.PinHash),
	)
	if err != nil {
		return err
	}

	var e string
	if p.LastError != nil {
		e = p.LastError.Error()
	}

	res := &app.PinbasePin{
		Hash:       string(p.ID),
		Aliases:    p.Aliases,
		WantPinned: p.WantPinned,
		Status:     p.Status.String(),
		LastError:  e,
	}

	// PinController_Show: end_implement
	return ctx.OK(res)
}

// Update runs the update action.
func (c *PinController) Update(ctx *app.UpdatePinContext) error {
	// PinController_Update: start_implement

	ps := c.P.PinService()

	err := ps.UpdatePin(
		pinbase.Hash(ctx.PartyHash),
		pinbase.Hash(ctx.PinHash),
		&pinbase.PinEdit{
			Aliases:    ctx.Payload.Aliases,
			WantPinned: *ctx.Payload.WantPinned,
		},
	)
	if err != nil {
		return err
	}

	p, err := ps.Pin(
		pinbase.Hash(ctx.PartyHash),
		pinbase.Hash(ctx.PinHash),
	)
	if err != nil {
		return err
	}

	var e string
	if p.LastError != nil {
		e = p.LastError.Error()
	}

	res := &app.PinbasePin{
		Hash:       string(p.ID),
		Aliases:    p.Aliases,
		WantPinned: p.WantPinned,
		Status:     p.Status.String(),
		LastError:  e,
	}

	// PinController_Update: end_implement
	return ctx.OK(res)
}
