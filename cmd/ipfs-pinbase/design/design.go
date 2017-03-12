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
			PartyHashParam()
		})
		Response(OK, PartyMedia)
		Response(NotFound)
	})

	Action("create", func() {
		Description("Create a party")
		Routing(POST(""))
		Payload(PartyCreatePayload, func() {
			Required("hash", "description")
		})
		Response(Created, "/parties/.+")
		Response(BadRequest, ErrorMedia)
	})

	Action("update", func() {
		Description("Change a party's description")
		Routing(PATCH("/:partyHash"))
		Params(func() {
			PartyHashParam()
		})
		Payload(PartyUpdatePayload)
		Response(OK, PartyMedia)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})

	Action("delete", func() {
		Description("Delete a party")
		Routing(DELETE("/:partyHash"))
		Params(func() {
			PartyHashParam()
		})
		Response(NoContent)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})
})

func PartyHashParam() {
	Param("partyHash", String, "Party Hash")
}

func PartyHash() {
	Attribute("hash", String, "The hash of the object describing the party")
}

func PartyDescription() {
	Attribute("description", String, "A helpful description of the party")
}

var PartyCreatePayload = Type("party-create-payload", func() {
	PartyHash()
	PartyDescription()
})

var PartyUpdatePayload = Type("party-update-payload", func() {
	PartyDescription()
})

var PartyMedia = MediaType("application/vnd.pinbase.party+json", func() {
	Description("A Pinbase Party")
	Attributes(func() {
		PartyHash()
		PartyDescription()
		Required("hash", "description")
	})
	View("default", func() {
		PartyHash()
		PartyDescription()
	})
})

var _ = Resource("pin", func() {
	Description("A thing to pin in IPFS")
	BasePath("/parties/:partyHash/pins")

	Action("list", func() {
		Description("List the pins under the party")
		Routing(GET(""))
		Params(func() {
			PartyHashParam()
		})
		Response(OK, func() {
			Media(CollectionOf(PinMedia))
		})
	})

	Action("show", func() {
		Description("Get the pin under the party by hash")
		Routing(GET("/:pinHash"))
		Params(func() {
			PartyHashParam()
			PinHashParam()
		})
		Response(OK, PinMedia)
		Response(NotFound)
	})

	Action("create", func() {
		Description("Create a pin under the party")
		Routing(POST(""))
		Params(func() {
			PartyHashParam()
		})
		Payload(PinCreatePayload, func() {
			Required("hash", "aliases", "want-pinned")
		})
		Response(Created, "/parties/.+/pins/.+")
		Response(BadRequest, ErrorMedia)
	})

	Action("update", func() {
		Description("Update a pin under the party")
		Routing(PATCH("/:pinHash"))
		Params(func() {
			PartyHashParam()
			PinHashParam()
		})
		Payload(PinUpdatePayload)
		Response(OK, PinMedia)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})

	Action("delete", func() {
		Description("Delete a pin under the party")
		Routing(DELETE("/:pinHash"))
		Params(func() {
			PartyHashParam()
			PinHashParam()
		})
		Response(NoContent)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})
})

func PinHashParam() {
	Param("pinHash", String, "Pin Hash")
}

func PinHash() {
	Attribute("hash", String, "The hash of the object to be pinned")
}

func PinAliases() {
	Attribute("aliases", ArrayOf(String), "Aliases for the pinned object")
}

func PinWantPinned() {
	Attribute("want-pinned", Boolean, "Indicates that the party wants to actually pin the object")
}

var PinCreatePayload = Type("pin-create-payload", func() {
	PinHash()
	PinAliases()
	PinWantPinned()
})

var PinUpdatePayload = Type("pin-update-payload", func() {
	PinAliases()
	PinWantPinned()
})

var PinMedia = MediaType("application/vnd.pinbase.pin+json", func() {
	Description("A Pin for a Party")
	Attributes(func() {
		PinHash()
		PinAliases()
		PinWantPinned()
		Attribute("status", String, "The status of the pin")
		Attribute("last-error", String, "Last pin error message")
		Required("hash", "aliases", "want-pinned", "status", "last-error")
	})
	View("default", func() {
		PinHash()
		PinAliases()
		PinWantPinned()
		Attribute("status")
		Attribute("last-error")
	})
})
