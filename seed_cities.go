package main

func lisbon() []seedVenue {
	c := "Lisbon"
	return []seedVenue{
		venue(c, "Alfama", "Tasca do Fado Azul", "restaurant", "Portuguese", "🐟", 4.5, 1284, 2,
			"Grilled sardines, vinho verde and live fado drifting in from next door.", "Live fado", "Terrace", "Local favourite"),
		venue(c, "Cais do Sodré", "Mercado Maré", "restaurant", "Seafood", "🦐", 4.5, 2310, 3,
			"Buzzing seafood counter famed for garlic prawns and clams Bulhão Pato.", "Shareable plates", "Late kitchen"),
		venue(c, "Bairro Alto", "Prego & Pão", "restaurant", "Sandwiches", "🥪", 4.0, 865, 1,
			"Garlicky steak sandwiches and ice-cold beers before the night kicks off.", "Cheap eats", "Quick bite"),
		venue(c, "Príncipe Real", "Jardim Verde", "restaurant", "Vegetarian", "🥗", 4.5, 642, 2,
			"Plant-based Portuguese petiscos in a candlelit garden courtyard.", "Vegan options", "Garden", "Date night"),
		venue(c, "Chiado", "Miradouro Rooftop", "bar", "Rooftop bar", "🌇", 4.5, 3120, 3,
			"Sunset cocktails with views over the Tagus and the castle.", "Sunset views", "Cocktails"),
		venue(c, "Rossio", "Ginjinha do Largo", "bar", "Cherry liqueur", "🍒", 4.5, 1980, 1,
			"Tiny standing-room bar pouring ginjinha in little chocolate cups.", "Local tradition", "Standing room"),
		venue(c, "Cais do Sodré", "Pensão Rosa", "bar", "Cocktail bar", "🍸", 4.0, 1543, 2,
			"A former guesthouse turned velvet-lined cocktail den on Pink Street.", "Cocktails", "Pink Street"),
		venue(c, "Alcântara", "Docas Club", "club", "Techno", "🎧", 4.0, 980, 2,
			"Riverside warehouse with a punishing sound system until sunrise.", "Open till 6am", "Big room"),
		venue(c, "Bairro Alto", "Bairro Beat", "club", "House & Disco", "🪩", 4.0, 712, 2,
			"Sweaty, friendly disco club where the whole street ends up at 2am.", "Disco", "No dress code"),
		venue(c, "Santos", "Kizomba Lounge", "club", "African beats", "💃", 4.5, 455, 2,
			"Kizomba, afrobeats and semba: come for the dancing, stay for the vibe.", "Dance classes", "Live percussion"),
	}
}

func barcelona() []seedVenue {
	c := "Barcelona"
	return []seedVenue{
		venue(c, "El Born", "La Bodega del Born", "restaurant", "Tapas", "🥘", 4.5, 3420, 2,
			"Elbow-to-elbow tapas bar: patatas bravas, croquetas and vermut on tap.", "Tapas", "Vermut", "Lively"),
		venue(c, "Barceloneta", "Can Marítim", "restaurant", "Paella & seafood", "🍤", 4.0, 2788, 3,
			"Seafront terrace serving black rice and paella for the table.", "Sea view", "Paella for two"),
		venue(c, "Gràcia", "Fum de Gràcia", "restaurant", "Catalan grill", "🔥", 4.5, 1105, 2,
			"Wood-fired calçots, butifarra and grilled vegetables on a leafy plaça.", "Outdoor plaça", "Wood-fired"),
		venue(c, "Eixample", "Sakura Eixample", "restaurant", "Japanese-Catalan", "🍣", 4.5, 864, 3,
			"Omakase counter mixing Catalan produce with Japanese technique.", "Chef's counter", "Date night"),
		venue(c, "Gothic Quarter", "El Vermut del Gòtic", "bar", "Vermouth bar", "🍷", 4.5, 1650, 1,
			"Barrel-lined bodega pouring house vermouth with olives and anchovies.", "Old-school", "Cheap drinks"),
		venue(c, "Eixample", "Terrat 360", "bar", "Rooftop bar", "🌆", 4.0, 2210, 3,
			"Rooftop pool bar with the Sagrada Família on the skyline.", "Rooftop", "Pool", "Views"),
		venue(c, "El Raval", "Absenta Raval", "bar", "Absinthe bar", "🥃", 4.0, 1032, 2,
			"Bohemian bar with absinthe fountains and art-covered walls.", "Quirky", "Absinthe"),
		venue(c, "Poblenou", "Fàbrica Nit", "club", "Indie & Pop", "🎸", 4.0, 3890, 2,
			"Five rooms of indie, pop and techno in a converted textile factory.", "Five rooms", "Open till 6am"),
		venue(c, "Port Olímpic", "Platja Club", "club", "Reggaeton & Charts", "🌴", 3.5, 2104, 3,
			"Beachfront club with a terrace that spills onto the sand.", "Beachfront", "Dress to impress"),
		venue(c, "Poble-sec", "Sala Ballroom", "club", "Electronic", "🎛️", 4.5, 1320, 2,
			"Historic 1900s ballroom hosting the city's best electronic nights.", "Historic venue", "Top DJs"),
	}
}

func amsterdam() []seedVenue {
	c := "Amsterdam"
	return []seedVenue{
		venue(c, "Jordaan", "De Gouden Pan", "restaurant", "Dutch", "🍲", 4.5, 1980, 2,
			"Dutch comfort food: stamppot, bitterballen and local beer by the canal.", "Canal side", "Cosy"),
		venue(c, "Centrum", "Rijsttafel Indah", "restaurant", "Indonesian", "🍛", 4.5, 2640, 3,
			"A feast of 18 small Indonesian dishes — the classic Amsterdam rijsttafel.", "Great for groups", "Feast"),
		venue(c, "De Pijp", "Pijp Pizza Lab", "restaurant", "Pizza", "🍕", 4.0, 1215, 1,
			"Blistered sourdough pizzas and natural wine in a buzzy little room.", "Cheap eats", "Natural wine"),
		venue(c, "Noord", "Noord Kantine", "restaurant", "Modern European", "🍽️", 4.0, 760, 2,
			"Converted shipyard canteen across the IJ, reached by the free ferry.", "Industrial", "Waterfront"),
		venue(c, "Centrum", "Proeflokaal Het Vat", "bar", "Jenever tasting", "🥃", 4.5, 1870, 1,
			"300-year-old tasting room: bend down for the first sip of jenever.", "Historic", "Local tradition"),
		venue(c, "Oost", "Brouwerij Oost", "bar", "Craft beer", "🍺", 4.5, 2430, 2,
			"Brewery taproom beside a windmill with a big summer terrace.", "Brewery", "Terrace"),
		venue(c, "Jordaan", "Gracht Cocktails", "bar", "Speakeasy", "🍸", 4.5, 690, 3,
			"Hidden behind a bookcase: candlelit speakeasy with inventive cocktails.", "Hidden entrance", "Cocktails"),
		venue(c, "Noord", "Loods Noord", "club", "Techno", "🎧", 4.5, 1590, 2,
			"24-hour licensed warehouse club with a legendary sound system.", "24h license", "Big room"),
		venue(c, "Leidseplein", "Leidse Lights", "club", "Party & Charts", "🪩", 3.5, 2875, 2,
			"Singalong anthems and cheap shots right on the main party square.", "Central", "No dress code"),
		venue(c, "Centrum", "Kerk Late", "club", "Live music & DJs", "🎤", 4.5, 4120, 2,
			"Former church turned concert hall, with club nights after the bands.", "Live bands", "Iconic venue"),
	}
}
