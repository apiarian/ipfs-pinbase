// Code generated by goagen v1.1.0-dirty, command line:
// $ goagen
// --design=github.com/apiarian/ipfs-pinbase/design
// --out=$(GOPATH)/src/github.com/apiarian/ipfs-pinbase
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

// ShowNodeContext provides the node show action context.
type ShowNodeContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	NodeHash string
}

// NewShowNodeContext parses the incoming request URL and body, performs validations and creates the
// context used by the node controller show action.
func NewShowNodeContext(ctx context.Context, service *goa.Service) (*ShowNodeContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	rctx := ShowNodeContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramNodeHash := req.Params["nodeHash"]
	if len(paramNodeHash) > 0 {
		rawNodeHash := paramNodeHash[0]
		rctx.NodeHash = rawNodeHash
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *ShowNodeContext) OK(r *PinbaseNode) error {
	ctx.ResponseData.Header().Set("Content-Type", "application/vnd.pinbase.node+json")
	return ctx.ResponseData.Service.Send(ctx.Context, 200, r)
}

// NotFound sends a HTTP response with status code 404.
func (ctx *ShowNodeContext) NotFound() error {
	ctx.ResponseData.WriteHeader(404)
	return nil
}
