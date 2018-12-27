// AUTO GENERATED FILE (by membufc proto compiler v0.0.21)
package gossiptopics

import (
	"context"
	"fmt"
	"github.com/orbs-network/orbs-spec/types/go/protocol/gossipmessages"
)

/////////////////////////////////////////////////////////////////////////////
// service LeanHelix

type LeanHelix interface {
	SendLeanHelixMessage(ctx context.Context, input *LeanHelixInput) (*EmptyOutput, error)
	RegisterLeanHelixHandler(handler LeanHelixHandler)
}

/////////////////////////////////////////////////////////////////////////////
// service LeanHelixHandler

type LeanHelixHandler interface {
	HandleLeanHelixMessage(ctx context.Context, input *LeanHelixInput) (*EmptyOutput, error)
}

/////////////////////////////////////////////////////////////////////////////
// message LeanHelixInput (non serializable)

type LeanHelixInput struct {
	RecipientsList *RecipientsList
	Message        *gossipmessages.LeanHelixMessage
}

func (x *LeanHelixInput) String() string {
	if x == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{RecipientsList:%s,Message:%s,}", x.StringRecipientsList(), x.StringMessage())
}

func (x *LeanHelixInput) StringRecipientsList() (res string) {
	res = x.RecipientsList.String()
	return
}

func (x *LeanHelixInput) StringMessage() (res string) {
	res = x.Message.String()
	return
}
