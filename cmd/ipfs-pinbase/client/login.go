// Code generated by goagen v1.1.0-dirty, command line:
// $ goagen
// --design=github.com/apiarian/ipfs-pinbase/cmd/ipfs-pinbase/design
// --out=$(GOPATH)/src/github.com/apiarian/ipfs-pinbase/cmd/ipfs-pinbase
// --version=v1.1.0-dirty
//
// API "pinbase": login Resource Client
//
// The content of this file is auto-generated, DO NOT MODIFY

package client

import (
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
)

// LoginLoginPath computes a request path to the login action of login.
func LoginLoginPath() string {

	return fmt.Sprintf("/login")
}

// Get a new JWT token
func (c *Client) LoginLogin(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewLoginLoginRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewLoginLoginRequest create the request corresponding to the login action endpoint of the login resource.
func (c *Client) NewLoginLoginRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.LoginBasicAuthSigner != nil {
		c.LoginBasicAuthSigner.Sign(req)
	}
	return req, nil
}
