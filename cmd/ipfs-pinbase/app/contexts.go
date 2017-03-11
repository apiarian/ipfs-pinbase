// Code generated by goagen v1.1.0-dirty, command line:
// $ goagen
// --design=github.com/apiarian/ipfs-pinbase/cmd/ipfs-pinbase/design
// --out=$(GOPATH)/src/github.com/apiarian/ipfs-pinbase/cmd/ipfs-pinbase
// --version=v1.1.0-dirty
//
// API "pinbase": Application Contexts
//
// The content of this file is auto-generated, DO NOT MODIFY

package app

import (
	"github.com/goadesign/goa"
	"golang.org/x/net/context"
)

// CreatePartyContext provides the party create action context.
type CreatePartyContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	Payload *CreatePartyPayload
}

// NewCreatePartyContext parses the incoming request URL and body, performs validations and creates the
// context used by the party controller create action.
func NewCreatePartyContext(ctx context.Context, service *goa.Service) (*CreatePartyContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := CreatePartyContext{Context: ctx, ResponseData: resp, RequestData: req}
	return &rctx, err
}

// createPartyPayload is the party create action payload.
type createPartyPayload struct {
	// A helpful description of the party
	Description *string `form:"description,omitempty" json:"description,omitempty" xml:"description,omitempty"`
	// The hash of the object describing the party
	Hash *string `form:"hash,omitempty" json:"hash,omitempty" xml:"hash,omitempty"`
}

// Validate runs the validation rules defined in the design.
func (payload *createPartyPayload) Validate() (err error) {
	if payload.Hash == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "hash"))
	}
	if payload.Description == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "description"))
	}
	return
}

// Publicize creates CreatePartyPayload from createPartyPayload
func (payload *createPartyPayload) Publicize() *CreatePartyPayload {
	var pub CreatePartyPayload
	if payload.Description != nil {
		pub.Description = *payload.Description
	}
	if payload.Hash != nil {
		pub.Hash = *payload.Hash
	}
	return &pub
}

// CreatePartyPayload is the party create action payload.
type CreatePartyPayload struct {
	// A helpful description of the party
	Description string `form:"description" json:"description" xml:"description"`
	// The hash of the object describing the party
	Hash string `form:"hash" json:"hash" xml:"hash"`
}

// Validate runs the validation rules defined in the design.
func (payload *CreatePartyPayload) Validate() (err error) {
	if payload.Hash == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "hash"))
	}
	if payload.Description == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`raw`, "description"))
	}
	return
}

// Created sends a HTTP response with status code 201.
func (ctx *CreatePartyContext) Created() error {
	ctx.ResponseData.WriteHeader(201)
	return nil
}

// BadRequest sends a HTTP response with status code 400.
func (ctx *CreatePartyContext) BadRequest(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 400, r)
}

// DeletePartyContext provides the party delete action context.
type DeletePartyContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	PartyHash string
}

// NewDeletePartyContext parses the incoming request URL and body, performs validations and creates the
// context used by the party controller delete action.
func NewDeletePartyContext(ctx context.Context, service *goa.Service) (*DeletePartyContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := DeletePartyContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramPartyHash := req.Params["partyHash"]
	if len(paramPartyHash) > 0 {
		rawPartyHash := paramPartyHash[0]
		rctx.PartyHash = rawPartyHash
	}
	return &rctx, err
}

// NoContent sends a HTTP response with status code 204.
func (ctx *DeletePartyContext) NoContent() error {
	ctx.ResponseData.WriteHeader(204)
	return nil
}

// BadRequest sends a HTTP response with status code 400.
func (ctx *DeletePartyContext) BadRequest(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 400, r)
}

// NotFound sends a HTTP response with status code 404.
func (ctx *DeletePartyContext) NotFound() error {
	ctx.ResponseData.WriteHeader(404)
	return nil
}

// ListPartyContext provides the party list action context.
type ListPartyContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
}

// NewListPartyContext parses the incoming request URL and body, performs validations and creates the
// context used by the party controller list action.
func NewListPartyContext(ctx context.Context, service *goa.Service) (*ListPartyContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := ListPartyContext{Context: ctx, ResponseData: resp, RequestData: req}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *ListPartyContext) OK(r PinbasePartyCollection) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.pinbase.party+json; type=collection")
	if r == nil {
		r = PinbasePartyCollection{}
	}
	return ctx.ResponseData.Service.Send(ctx.Context, 200, r)
}

// ShowPartyContext provides the party show action context.
type ShowPartyContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	PartyHash string
}

// NewShowPartyContext parses the incoming request URL and body, performs validations and creates the
// context used by the party controller show action.
func NewShowPartyContext(ctx context.Context, service *goa.Service) (*ShowPartyContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := ShowPartyContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramPartyHash := req.Params["partyHash"]
	if len(paramPartyHash) > 0 {
		rawPartyHash := paramPartyHash[0]
		rctx.PartyHash = rawPartyHash
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *ShowPartyContext) OK(r *PinbaseParty) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.pinbase.party+json")
	return ctx.ResponseData.Service.Send(ctx.Context, 200, r)
}

// NotFound sends a HTTP response with status code 404.
func (ctx *ShowPartyContext) NotFound() error {
	ctx.ResponseData.WriteHeader(404)
	return nil
}

// UpdatePartyContext provides the party update action context.
type UpdatePartyContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	PartyHash string
	Payload   *PartyUpdatePayload
}

// NewUpdatePartyContext parses the incoming request URL and body, performs validations and creates the
// context used by the party controller update action.
func NewUpdatePartyContext(ctx context.Context, service *goa.Service) (*UpdatePartyContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := UpdatePartyContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramPartyHash := req.Params["partyHash"]
	if len(paramPartyHash) > 0 {
		rawPartyHash := paramPartyHash[0]
		rctx.PartyHash = rawPartyHash
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *UpdatePartyContext) OK(r *PinbaseParty) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.pinbase.party+json")
	return ctx.ResponseData.Service.Send(ctx.Context, 200, r)
}

// BadRequest sends a HTTP response with status code 400.
func (ctx *UpdatePartyContext) BadRequest(r error) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.goa.error")
	return ctx.ResponseData.Service.Send(ctx.Context, 400, r)
}

// NotFound sends a HTTP response with status code 404.
func (ctx *UpdatePartyContext) NotFound() error {
	ctx.ResponseData.WriteHeader(404)
	return nil
}
