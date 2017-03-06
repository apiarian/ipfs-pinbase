// Code generated by goagen v1.1.0-dirty, command line:
// $ goagen
// --design=github.com/apiarian/ipfs-pinbase/cmd/ipfs-pinbase/design
// --out=$(GOPATH)/src/github.com/apiarian/ipfs-pinbase/cmd/ipfs-pinbase
// --version=v1.1.0-dirty
//
// API "pinbase": Application Media Types
//
// The content of this file is auto-generated, DO NOT MODIFY

package client

import (
	"github.com/goadesign/goa"
	"net/http"
)

// DecodeErrorResponse decodes the ErrorResponse instance encoded in resp body.
func (c *Client) DecodeErrorResponse(resp *http.Response) (*goa.ErrorResponse, error) {
	var decoded goa.ErrorResponse
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return &decoded, err
}

// An IPFS node (default view)
//
// Identifier: application/vnd.pinbase.node+json; view=default
type PinbaseNode struct {
	// The API URL for the node, possibly relative to the pinbase (i.e. localhost)
	APIURL string `form:"api-url" json:"api-url" xml:"api-url"`
	// A helpful description of the node
	Description string `form:"description" json:"description" xml:"description"`
	// The nodes' unique hash
	Hash string `form:"hash" json:"hash" xml:"hash"`
}

// Validate validates the PinbaseNode media type instance.
func (mt *PinbaseNode) Validate() (err error) {
	if mt.Hash == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "hash"))
	}
	if mt.Description == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "description"))
	}
	if mt.APIURL == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "api-url"))
	}
	return
}

// DecodePinbaseNode decodes the PinbaseNode instance encoded in resp body.
func (c *Client) DecodePinbaseNode(resp *http.Response) (*PinbaseNode, error) {
	var decoded PinbaseNode
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return &decoded, err
}

// PinbaseNodeCollection is the media type for an array of PinbaseNode (default view)
//
// Identifier: application/vnd.pinbase.node+json; type=collection; view=default
type PinbaseNodeCollection []*PinbaseNode

// Validate validates the PinbaseNodeCollection media type instance.
func (mt PinbaseNodeCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// DecodePinbaseNodeCollection decodes the PinbaseNodeCollection instance encoded in resp body.
func (c *Client) DecodePinbaseNodeCollection(resp *http.Response) (PinbaseNodeCollection, error) {
	var decoded PinbaseNodeCollection
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return decoded, err
}