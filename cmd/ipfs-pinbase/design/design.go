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
	BasePath("/api")

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

var _ = Resource("party", func() {
	Description("The Pinbase Party resource")
	BasePath("/parties")

	Action("list", func() {
		Description("List the parties available in this pinbase")
		Routing(GET(""))
		Response(OK, func() {
			Media(CollectionOf(PartyMedia))
		})
	})

	Action("show", func() {
		Description("Get the party by hash")
		Routing(GET("/:partyHash"))
		Params(func() {
			Param("partyHash", String, "Party Hash")
		})
		Response(OK, PartyMedia)
		Response(NotFound)
	})

	Action("create", func() {
		Description("Create a party")
		Routing(POST(""))
		Payload(PartyPayload, func() {
			Required("hash", "description")
		})
		Response(Created, "/parties/.+")
		Response(BadRequest, ErrorMedia)
	})

	Action("update", func() {
		Description("Change a party's description")
		Routing(PATCH("/:partyHash"))
		Params(func() {
			Param("partyHash", String, "Party Hash")
		})
		Payload(PartyPayload)
		Response(OK, PartyMedia)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})

	Action("delete", func() {
		Description("Delete a party")
		Routing(DELETE("/:partyHash"))
		Params(func() {
			Param("partyHash", String, "Party Hash")
		})
		Response(NoContent)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})
})

var PartyPayload = Type("party-payload", func() {
	Attribute("hash", String, "The hash of the object describing the party")
	Attribute("description", String, "A helpful description of the party")
})

var PartyMedia = MediaType("application/vnd.pinbase.party+json", func() {
	Description("A Pinbase Party")
	Reference(PartyPayload)
	Attributes(func() {
		Attribute("hash")
		Attribute("description")
		Required("hash", "description")
	})
	View("default", func() {
		Attribute("hash")
		Attribute("description")
	})
})
