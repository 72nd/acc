package api

import (
	"github.com/72nd/acc/pkg/schema"
)

// fromAccParties converts a schema.Parties collection into the API representation.
func fromAccParties(parties []schema.Party) Parties {
	pty := make(Parties, len(parties))
	for i := range parties {
		pty[i] = fromAccParty(parties[i])
	}
	return pty
}

// fromAccParty converts a schema.Party object into the API representation.
func fromAccParty(party schema.Party) Party {
	partyType := new(int)
	*partyType = int(party.PartyType)
	postalCode := new(int)
	*postalCode = party.PostalCode
	streetNr := new(int)
	*streetNr = party.StreetNr

	return Party{
		Id:         &party.Id,
		Identifier: &party.Identifier,
		Name:       &party.Name,
		PartyType:  partyType,
		Place:      &party.Place,
		PostalCode: postalCode,
		Street:     &party.Street,
		StreetNr:   streetNr,
	}
}
