package synchronization_test

import (
	"context"
	"github.com/orbs-network/orbs-network-go/instrumentation/log"
	"github.com/orbs-network/orbs-network-go/synchronization"
	"github.com/stretchr/testify/require"
	"sync/atomic"
	"testing"
	"time"
)

type report struct {
	message string
	fields  []*log.Field
}

type collector struct {
	errors chan report
}

func (c *collector) Error(message string, fields ...*log.Field) {
	c.errors <- report{message, fields}
}

func mockLogger() *collector {
	c := &collector{errors: make(chan report)}
	return c
}

func getExpected(startTime, endTime time.Time, tickTime time.Duration) uint32 {
	duration := endTime.Sub(startTime)
	expected := uint32((duration.Seconds() * 1000) / (tickTime.Seconds() * 1000))
	return expected
}

func TestPeriodicalTriggerStartsOk(t *testing.T) {
	logger := mockLogger()
	var x uint32
	start := time.Now()
	tickTime := 5 * time.Millisecond
	p := synchronization.NewPeriodicalTrigger(context.Background(), tickTime, logger, func() { atomic.AddUint32(&x, 1) }, nil)
	time.Sleep(time.Millisecond * 30)
	expected := getExpected(start, time.Now(), tickTime)
	require.True(t, expected/2 < atomic.LoadUint32(&x), "expected more than %d ticks, but got %d", expected/2, atomic.LoadUint32(&x))
	p.Stop()
}

func TestTriggerInternalMetrics(t *testing.T) {
	logger := mockLogger()
	var x uint32
	start := time.Now()
	tickTime := 5 * time.Millisecond
	p := synchronization.NewPeriodicalTrigger(context.Background(), tickTime, logger, func() { atomic.AddUint32(&x, 1) }, nil)
	time.Sleep(time.Millisecond * 30)
	expected := getExpected(start, time.Now(), tickTime)
	require.True(t, expected/2 < atomic.LoadUint32(&x), "expected more than %d ticks, but got %d", expected/2, atomic.LoadUint32(&x))
	require.True(t, uint64(expected/2) < p.TimesTriggered(), "expected more than %d ticks, but got %d (metric)", expected/2, p.TimesTriggered())
	p.Stop()
}

func TestPeriodicalTrigger_Stop(t *testing.T) {
	logger := mockLogger()
	x := 0
	p := synchronization.NewPeriodicalTrigger(context.Background(), time.Millisecond*2, logger, func() { x++ }, nil)
	p.Stop()
	time.Sleep(3 * time.Millisecond)
	require.Equal(t, 0, x, "expected no ticks")
}

func TestPeriodicalTrigger_StopAfterTrigger(t *testing.T) {
	logger := mockLogger()
	x := 0
	p := synchronization.NewPeriodicalTrigger(context.Background(), time.Millisecond, logger, func() { x++ }, nil)
	time.Sleep(time.Microsecond * 1100)
	p.Stop()
	xValueOnStop := x
	time.Sleep(time.Millisecond * 5)
	require.Equal(t, xValueOnStop, x, "expected one tick due to stop")
}

func TestPeriodicalTriggerStopOnContextCancel(t *testing.T) {
	logger := mockLogger()
	ctx, cancel := context.WithCancel(context.Background())
	x := 0
	synchronization.NewPeriodicalTrigger(ctx, time.Millisecond*2, logger, func() { x++ }, nil)
	cancel()
	time.Sleep(3 * time.Millisecond)
	require.Equal(t, 0, x, "expected no ticks")
}

func TestPeriodicalTriggerStopWorksWhenContextIsCancelled(t *testing.T) {
	logger := mockLogger()
	ctx, cancel := context.WithCancel(context.Background())
	x := 0
	p := synchronization.NewPeriodicalTrigger(ctx, time.Millisecond*2, logger, func() { x++ }, nil)
	cancel()
	time.Sleep(3 * time.Millisecond)
	require.Equal(t, 0, x, "expected no ticks")
	p.Stop()
	require.Equal(t, 0, x, "expected stop to not block")
}

func TestPeriodicalTriggerStopOnContextCancelWithStopAction(t *testing.T) {
	logger := mockLogger()
	ctx, cancel := context.WithCancel(context.Background())
	x := 0
	synchronization.NewPeriodicalTrigger(ctx, time.Millisecond*2, logger, func() { x++ }, func() { x = 20 })
	cancel()
	time.Sleep(time.Millisecond) // yield
	require.Equal(t, 20, x, "expected x to have the stop value")
}

func TestPeriodicalTriggerRunsOnStopAction(t *testing.T) {
	logger := mockLogger()
	latch := make(chan struct{})
	x := 0
	p := synchronization.NewPeriodicalTrigger(context.Background(),
		time.Second,
		logger,
		func() { x++ },
		func() {
			x = 20
			latch <- struct{}{}
		})
	p.Stop()
	<-latch // wait for stop to happen...
	require.Equal(t, 20, x, "expected x to have the stop value")
}

func TestPeriodicalTriggerKeepsGoingOnPanic(t *testing.T) {
	t.Skip("LongLived is broken, try again when its fixed")
	logger := mockLogger()
	x := 0
	p := synchronization.NewPeriodicalTrigger(context.Background(),
		time.Millisecond,
		logger,
		func() {
			x++
			panic("we should not see this other than the logs")
		},
		nil)
	time.Sleep(5 * time.Millisecond)
	p.Stop()
	require.True(t, x > 1, "expected trigger to have ticked more than once (even though it panics) %d", x)
}
