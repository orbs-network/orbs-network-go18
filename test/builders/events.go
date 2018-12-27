package builders

import (
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
)

func PackedEventsArrayEncode(eventBuilders []*protocol.EventBuilder) primitives.PackedEventsArray {
	eventsArray := (&protocol.EventsArrayBuilder{Events: eventBuilders}).Build()
	return eventsArray.RawEventsArray()
}
