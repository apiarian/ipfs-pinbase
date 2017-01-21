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
})

var _ = Resource("node", func() {
	BasePath("/nodes")
	DefaultMedia(NodeMedia)

	Action("show", func() {
		Description("Get node by hash")
		Routing(GET("/:nodeHash"))
		Params(func() {
			Param("nodeHash", String, "Node Hash")
		})
		Response(OK)
		Response(NotFound)
	})
})

var NodeMedia = MediaType("application/vnd.pinbase.node+json", func() {
	Description("An IPFS node")
	Attributes(func() {
		Attribute("hash", String, "The nodes' unique hash")
		Attribute("description", String, "A helpful description of the node")
		Attribute("api-url", String, "The API URL for the node, possibly relative to the pinbase (i.e. localhost)")
		Required("hash", "description", "api-url")
	})
	View("default", func() {
		Attribute("hash")
		Attribute("description")
		Attribute("api-url")
	})
})
