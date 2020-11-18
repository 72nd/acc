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

// fromApiParty converts a API Party into a schema.Party object. If no identifier is given
// it will used the nextIdent value. An error is returned if some input data isn't valid.
func fromApiParty(party Party, nextIdentifier string) (schema.Party, error) {
	rsl := schema.NewPartyWithUuid()
	rsl.Name = stringValueOrDefault(party.Name)
	if party.Identifier != nil {
		rsl.Identifier = *party.Identifier
	} else {
		rsl.Identifier = nextIdentifier
	}
	rsl.StreetNr = intValueOrDefault(party.StreetNr)
	rsl.PostalCode = intValueOrDefault(party.PostalCode)
	rsl.Place = stringValueOrDefault(party.Place)
	return rsl, nil
}

// stringValueOrDefault takes a string reference and returns the content as a value or an
// empty string if the pointer is nil.
func stringValueOrDefault(ele *string) string {
	if ele != nil {
		return *ele
	}
	return ""
}

// intValueOrDefault takes a int reference and returns the content as a value or an
// empty string if the pointer is nil.
func intValueOrDefault(ele *int) int {
	if ele != nil {
		return *ele
	}
	return 0
}
