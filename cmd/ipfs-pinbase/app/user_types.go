// Code generated by goagen v1.1.0-dirty, command line:
// $ goagen
// --design=github.com/apiarian/ipfs-pinbase/cmd/ipfs-pinbase/design
// --out=$(GOPATH)/src/github.com/apiarian/ipfs-pinbase/cmd/ipfs-pinbase
// --version=v1.1.0-dirty
//
// API "pinbase": Application User Types
//
// The content of this file is auto-generated, DO NOT MODIFY

package app

// nodePayload user type.
type nodePayload struct {
	// The API URL for the node, possibly relative to the pinbase (i.e. localhost)
	APIURL *string `form:"api-url,omitempty" json:"api-url,omitempty" xml:"api-url,omitempty"`
	// A helpful description of the node
	Description *string `form:"description,omitempty" json:"description,omitempty" xml:"description,omitempty"`
}

// Publicize creates NodePayload from nodePayload
func (ut *nodePayload) Publicize() *NodePayload {
	var pub NodePayload
	if ut.APIURL != nil {
		pub.APIURL = ut.APIURL
	}
	if ut.Description != nil {
		pub.Description = ut.Description
	}
	return &pub
}

// NodePayload user type.
type NodePayload struct {
	// The API URL for the node, possibly relative to the pinbase (i.e. localhost)
	APIURL *string `form:"api-url,omitempty" json:"api-url,omitempty" xml:"api-url,omitempty"`
	// A helpful description of the node
	Description *string `form:"description,omitempty" json:"description,omitempty" xml:"description,omitempty"`
}
