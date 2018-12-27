package synchronization

import (
	"context"
	"github.com/orbs-network/orbs-network-go/instrumentation/log"
	"github.com/orbs-network/orbs-network-go/test"
	"github.com/stretchr/testify/require"
	"sync/atomic"
	"testing"
)

func TestWaitForBlockOutsideOfGraceFailsImmediately(t *testing.T) {
	test.WithContext(func(ctx context.Context) {
		tracker := NewBlockTracker(log.GetLogger(), 1, 1)

		err := tracker.WaitForBlock(ctx, 3)
		require.EqualError(t, err, "requested future block outside of grace range", "did not fail immediately")
	})
}

func TestWaitForBlockWithinGraceFailsWhenContextEnds(t *testing.T) {
	test.WithContext(func(parentCtx context.Context) {
		ctx, cancel := context.WithCancel(parentCtx)
		tracker := NewBlockTracker(log.GetLogger(), 1, 1)
		cancel()
		err := tracker.WaitForBlock(ctx, 2)
		require.EqualError(t, err, "aborted while waiting for block at height 2: context canceled", "did not fail as expected")
	})
}

func TestWaitForBlockWithinGraceDealsWithIntegerUnderflow(t *testing.T) {
	test.WithContext(func(parentCtx context.Context) {
		ctx, cancel := context.WithCancel(parentCtx)
		tracker := NewBlockTracker(log.GetLogger(), 0, 5)
		cancel()
		err := tracker.WaitForBlock(ctx, 2)
		require.EqualError(t, err, "aborted while waiting for block at height 2: context canceled", "did not fail as expected")
	})
}

func TestWaitForBlockWithinGraceReturnsWhenBlockHeightReachedBeforeContextEnds(t *testing.T) {
	test.WithContext(func(ctx context.Context) {
		tracker := NewBlockTracker(log.GetLogger(), 1, 2)

		var waitCount int32
		internalWaitChan := make(chan int32)
		tracker.fireOnWait = func() {
			internalWaitChan <- atomic.AddInt32(&waitCount, 1)
		}

		doneWait := make(chan error)
		go func() {
			doneWait <- tracker.WaitForBlock(ctx, 3)
		}()

		require.EqualValues(t, 1, <-internalWaitChan, "did not block before the first increment")
		tracker.IncrementHeight()
		require.EqualValues(t, 2, <-internalWaitChan, "did not block before the second increment")
		tracker.IncrementHeight()

		require.NoError(t, <-doneWait, "did not return as expected")
	})
}

func TestWaitForBlockWithinGraceSupportsTwoConcurrentWaiters(t *testing.T) {
	test.WithContext(func(ctx context.Context) {
		tracker := NewBlockTracker(log.GetLogger(), 1, 1)

		var waitCount int32
		internalWaitChan := make(chan int32)
		tracker.fireOnWait = func() {
			internalWaitChan <- atomic.AddInt32(&waitCount, 1)
		}

		doneWait := make(chan error)
		waiter := func() {
			doneWait <- tracker.WaitForBlock(ctx, 2)
		}
		go waiter()
		go waiter()

		selectIterationsBeforeIncrement := <-internalWaitChan
		require.EqualValues(t, 1, selectIterationsBeforeIncrement, "did not enter select before returning")
		selectIterationsBeforeIncrement = <-internalWaitChan
		require.EqualValues(t, 2, selectIterationsBeforeIncrement, "did not enter select before returning")

		tracker.IncrementHeight()

		require.NoError(t, <-doneWait, "first waiter did not return as expected")
		require.NoError(t, <-doneWait, "second waiter did not return as expected")
	})
}
