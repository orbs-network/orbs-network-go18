package adapter

import (
	"context"
	"github.com/orbs-network/orbs-network-go/services/gossip/adapter"
	"github.com/orbs-network/orbs-network-go/synchronization/supervised"
	"github.com/orbs-network/orbs-network-go/test"
	"runtime"
	"sync"
	"time"
)

type failingTamperer struct {
	predicate MessagePredicate
	transport *TamperingTransport
}

func (o *failingTamperer) maybeTamper(ctx context.Context, data *adapter.TransportData) (error, bool) {
	if o.predicate(data) {
		return &adapter.ErrTransportFailed{Data: data}, true
	}

	return nil, false
}

func (o *failingTamperer) Release(ctx context.Context) {
	o.transport.removeOngoingTamperer(o)
}

type duplicatingTamperer struct {
	predicate MessagePredicate
	transport *TamperingTransport
}

func (o *duplicatingTamperer) maybeTamper(ctx context.Context, data *adapter.TransportData) (error, bool) {
	if o.predicate(data) {
		supervised.GoOnce(o.transport.logger, func() {
			time.Sleep(10 * time.Millisecond)
			o.transport.sendToPeers(ctx, data)
		})
	}
	return nil, false
}

func (o *duplicatingTamperer) Release(ctx context.Context) {
	o.transport.removeOngoingTamperer(o)
}

type delayingTamperer struct {
	predicate MessagePredicate
	transport *TamperingTransport
	duration  func() time.Duration
}

func (o *delayingTamperer) maybeTamper(ctx context.Context, data *adapter.TransportData) (error, bool) {
	if o.predicate(data) {
		supervised.GoOnce(o.transport.logger, func() {
			time.Sleep(o.duration())
			o.transport.sendToPeers(ctx, data)
		})
		return nil, true
	}

	return nil, false
}

func (o *delayingTamperer) Release(ctx context.Context) {
	o.transport.removeOngoingTamperer(o)
}

type corruptingTamperer struct {
	predicate MessagePredicate
	transport *TamperingTransport
	ctrlRand  *test.ControlledRand
}

func (o *corruptingTamperer) maybeTamper(ctx context.Context, data *adapter.TransportData) (error, bool) {
	if o.predicate(data) {
		for i := 0; i < 10; i++ {
			x := o.ctrlRand.Intn(len(data.Payloads))
			y := o.ctrlRand.Intn(len(data.Payloads[x]))
			data.Payloads[x][y] = 0
		}
	}
	return nil, false
}

func (o *corruptingTamperer) Release(ctx context.Context) {
	o.transport.removeOngoingTamperer(o)
}

type pausingTamperer struct {
	predicate MessagePredicate
	transport *TamperingTransport
	messages  []*adapter.TransportData
	lock      *sync.Mutex
}

func (o *pausingTamperer) maybeTamper(ctx context.Context, data *adapter.TransportData) (error, bool) {
	if o.predicate(data) {
		o.lock.Lock()
		defer o.lock.Unlock()
		o.messages = append(o.messages, data)
		return nil, true
	}

	return nil, false
}

func (o *pausingTamperer) Release(ctx context.Context) {
	o.transport.removeOngoingTamperer(o)
	for _, message := range o.messages {
		o.transport.Send(ctx, message)
		runtime.Gosched() // TODO(v1): this is required or else messages arrive in the opposite order after resume (supposedly fixed now when we moved to channels in transport)
	}
}

type latchingTamperer struct {
	predicate MessagePredicate
	transport *TamperingTransport
	cond      *sync.Cond
}

func (o *latchingTamperer) Remove() {
	o.transport.removeLatchingTamperer(o)
}

func (o *latchingTamperer) Wait() {
	o.cond.Wait()
}
