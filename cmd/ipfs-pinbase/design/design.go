package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("pinbase", func() {
	Title("pinbase")
	Description("The IPFS-pinbase API")

	Version("0.1")

	Contact(func() {
		Name("Aleksandr Pasechnik")
		Email("al@megamicron.net")
		URL("https://megamicron.net")
	})

	License(func() {
		Name("MIT")
	})

	Host("localhost:3000")
	Scheme("http")
	BasePath("")

	Consumes("application/json")
	Produces("application/json")

	ResponseTemplate(Created, func(pattern string) {
		Description("Resource created")
		Status(201)
		Headers(func() {
			Header("Location", String, "href to the created resource", func() {
				Pattern(pattern)
			})
		})
	})
})

var LoginBasicAuth = BasicAuthSecurity("LoginBasicAuth")

var JWT = JWTSecurity("jwt", func() {
	Header("Authorization")
	TokenURL("/login")

	Scope("node:view")
	Scope("node:create")
	Scope("node:edit")
	Scope("node:delete")
})

var _ = Resource("login", func() {
	Description("The login resrouce to obtain a token")

	Action("login", func() {
		Description("Get a new JWT token")
		Routing(POST("/login"))
		Security(LoginBasicAuth)
		Response(NoContent, func() {
			Headers(func() {
				Header("Authorization", String, "The new JWT")
			})
		})
	})
})

var _ = Resource("node", func() {
	Description("The IPFS node resrouce")
	BasePath("/nodes")

	Action("list", func() {
		Description("List the nodes available to this pinbase")
		Routing(GET(""))
		Security(JWT, func() {
			Scope("node:view")
		})
		Response(OK, func() {
			Media(CollectionOf(NodeMedia))
		})
	})

	Action("show", func() {
		Description("Get node by hash")
		Routing(GET("/:nodeHash"))
		Params(func() {
			Param("nodeHash", String, "Node Hash")
		})
		Security(JWT, func() {
			Scope("node:view")
		})
		Response(OK, NodeMedia)
		Response(NotFound)
	})

	Action("create", func() {
		Description("Connect to a node")
		Routing(POST(""))
		Payload(NodePayload, func() {
			Required("api-url", "description")
		})
		Security(JWT, func() {
			Scope("node:create")
		})
		Response(Created, "/nodes/.+")
		Response(BadRequest, ErrorMedia)
	})

	Action("update", func() {
		Description("Change a node's address (must be the same node-id) or description")
		Routing(PATCH("/:nodeHash"))
		Params(func() {
			Param("nodeHash", String, "Node Hash")
		})
		Payload(NodePayload)
		Security(JWT, func() {
			Scope("node:edit")
		})
		Response(OK, NodeMedia)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})

	Action("delete", func() {
		Description("Delete a node")
		Routing(DELETE("/:nodeHash"))
		Params(func() {
			Param("nodeHash", String, "Node Hash")
		})
		Security(JWT, func() {
			Scope("node:delete")
		})
		Response(NoContent)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})
})

var NodePayload = Type("node-payload", func() {
	Attribute("api-url", String, "The API URL for the node, possibly relative to the pinbase (i.e. localhost)")
	Attribute("description", String, "A helpful description of the node")
})

var NodeMedia = MediaType("application/vnd.pinbase.node+json", func() {
	Description("An IPFS node")
	Reference(NodePayload)
	Attributes(func() {
		Attribute("hash", String, "The nodes' unique hash")
		Attribute("description")
		Attribute("api-url")
		Required("hash", "description", "api-url")
	})
	View("default", func() {
		Attribute("hash")
		Attribute("description")
		Attribute("api-url")
	})
})
