package adapter

import (
	"context"
	"fmt"
	"github.com/orbs-network/membuffers/go"
	"github.com/orbs-network/orbs-network-go/config"
	"github.com/orbs-network/orbs-network-go/instrumentation/log"
	"github.com/orbs-network/orbs-network-go/instrumentation/metric"
	"github.com/orbs-network/orbs-network-go/instrumentation/trace"
	"github.com/orbs-network/orbs-network-go/synchronization/supervised"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol/gossipmessages"
	"github.com/pkg/errors"
	"net"
	"sync"
	"time"
)

const MAX_PAYLOADS_IN_MESSAGE = 100000
const MAX_PAYLOAD_SIZE_BYTES = 10 * 1024 * 1024

var LogTag = log.String("adapter", "gossip")

type metrics struct {
	incomingConnectionAcceptErrors    *metric.Gauge
	incomingConnectionTransportErrors *metric.Gauge
	outgoingConnectionFailedSend      *metric.Gauge
	outgoingConnectionFailedKeepalive *metric.Gauge
}

type directTransport struct {
	config config.GossipTransportConfig
	logger log.BasicLogger

	outgoingPeerQueues map[string]chan *TransportData // does not require mutex to read

	mutex                       *sync.RWMutex
	transportListenerUnderMutex TransportListener
	serverListeningUnderMutex   bool
	serverPort                  int

	metrics *metrics
}

func getMetrics(registry metric.Registry) *metrics {
	return &metrics{
		incomingConnectionAcceptErrors:    registry.NewGauge("Gossip.IncomingConnection.AcceptErrors"),
		incomingConnectionTransportErrors: registry.NewGauge("Gossip.IncomingConnection.TransportErrors"),
		outgoingConnectionFailedSend:      registry.NewGauge("Gossip.OutgoingConnection.FailedSendErrors"),
		outgoingConnectionFailedKeepalive: registry.NewGauge("Gossip.OutgoingConnection.FailedKeepaliveErrors"),
	}
}

func NewDirectTransport(ctx context.Context, config config.GossipTransportConfig, logger log.BasicLogger, registry metric.Registry) Transport {
	t := &directTransport{
		config: config,
		logger: logger.WithTags(LogTag),

		outgoingPeerQueues: make(map[string]chan *TransportData),

		mutex:   &sync.RWMutex{},
		metrics: getMetrics(registry),
	}

	// client channels (not under mutex, before all goroutines)
	for peerNodeAddress := range t.config.GossipPeers(0) {
		if peerNodeAddress != t.config.NodeAddress().KeyForMap() {
			t.outgoingPeerQueues[peerNodeAddress] = make(chan *TransportData)
		}
	}

	// server goroutine
	supervised.GoForever(ctx, t.logger, func() {
		t.serverMainLoop(ctx, t.config.GossipListenPort())
	})

	// client goroutines
	for peerNodeAddress, peer := range t.config.GossipPeers(0) {
		if peerNodeAddress != t.config.NodeAddress().KeyForMap() {
			peerAddress := fmt.Sprintf("%s:%d", peer.GossipEndpoint(), peer.GossipPort())
			closureSafePeerNodeKey := peerNodeAddress
			supervised.GoForever(ctx, t.logger, func() {
				t.clientMainLoop(ctx, peerAddress, t.outgoingPeerQueues[closureSafePeerNodeKey])
			})
		}
	}

	return t
}

func (t *directTransport) RegisterListener(listener TransportListener, listenerNodeAddress primitives.NodeAddress) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.transportListenerUnderMutex = listener
}

// TODO(https://github.com/orbs-network/orbs-network-go/issues/182): we are not currently respecting any intents given in ctx (added in context refactor)
func (t *directTransport) Send(ctx context.Context, data *TransportData) error {
	switch data.RecipientMode {
	case gossipmessages.RECIPIENT_LIST_MODE_BROADCAST:
		for _, peerQueue := range t.outgoingPeerQueues {
			peerQueue <- data
		}
		return nil
	case gossipmessages.RECIPIENT_LIST_MODE_LIST:
		for _, recipientPublicKey := range data.RecipientNodeAddresses {
			if peerQueue, found := t.outgoingPeerQueues[recipientPublicKey.KeyForMap()]; found {
				peerQueue <- data
			} else {
				return errors.Errorf("unknown recipient public key: %s", recipientPublicKey.String())
			}
		}
		return nil
	case gossipmessages.RECIPIENT_LIST_MODE_ALL_BUT_LIST:
		panic("Not implemented")
	}
	return errors.Errorf("unknown recipient mode: %s", data.RecipientMode.String())
}

func (t *directTransport) serverListenForIncomingConnections(ctx context.Context, listenPort uint16) (net.Listener, error) {
	// TODO(v1): migrate to ListenConfig which has better support of contexts (go 1.11 required)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", listenPort))
	if err != nil {
		return nil, err
	}

	// this goroutine will shut down the server gracefully when context is done
	go func() {
		<-ctx.Done()
		t.mutex.Lock()
		defer t.mutex.Unlock()
		t.serverListeningUnderMutex = false
		listener.Close()
	}()

	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.serverListeningUnderMutex = true

	return listener, err
}

func (t *directTransport) isServerListening() bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.serverListeningUnderMutex
}

func (t *directTransport) serverMainLoop(parentCtx context.Context, listenPort uint16) {
	listener, err := t.serverListenForIncomingConnections(parentCtx, listenPort)
	if err != nil {
		err = errors.Wrapf(err, "gossip transport cannot listen on port %d", listenPort)
		t.logger.Error(err.Error())
		panic(err)
	}

	t.serverPort = listener.Addr().(*net.TCPAddr).Port
	t.logger.Info("gossip transport server listening", log.Uint32("port", uint32(t.serverPort)))

	for {
		if parentCtx.Err() != nil {
			t.logger.Info("ending server main loop (system shutting down)")
		}

		ctx := trace.NewContext(parentCtx, "Gossip.Transport.TCP.Server")

		conn, err := listener.Accept()
		if err != nil {
			if !t.isServerListening() {
				t.logger.Info("incoming connection accept stopped since server is shutting down", trace.LogFieldFrom(ctx))
				return
			}
			t.metrics.incomingConnectionAcceptErrors.Inc()
			t.logger.Info("incoming connection accept error", log.Error(err), trace.LogFieldFrom(ctx))
			continue
		}
		supervised.GoOnce(t.logger, func() {
			t.serverHandleIncomingConnection(ctx, conn)
		})
	}
}

func (t *directTransport) serverHandleIncomingConnection(ctx context.Context, conn net.Conn) {
	t.logger.Info("successful incoming gossip transport connection", log.String("peer", conn.RemoteAddr().String()), trace.LogFieldFrom(ctx))
	// TODO(https://github.com/orbs-network/orbs-network-go/issues/182): add a white list for IPs we're willing to accept connections from
	// TODO(https://github.com/orbs-network/orbs-network-go/issues/182): make sure each IP from the white list connects only once

	for {
		payloads, err := t.receiveTransportData(ctx, conn)
		if err != nil {
			t.metrics.incomingConnectionTransportErrors.Inc()
			t.logger.Info("failed receiving transport data, disconnecting", log.Error(err), log.String("peer", conn.RemoteAddr().String()), trace.LogFieldFrom(ctx))
			conn.Close()
			return
		}

		// notify if not keepalive
		if len(payloads) > 0 {
			t.notifyListener(ctx, payloads)
		}
	}
}

func (t *directTransport) receiveTransportData(ctx context.Context, conn net.Conn) ([][]byte, error) {
	// TODO(https://github.com/orbs-network/orbs-network-go/issues/182): think about timeout policy on receive, we might not want it
	timeout := t.config.GossipNetworkTimeout()
	res := [][]byte{}

	// receive num payloads
	sizeBuffer, err := readTotal(ctx, conn, 4, timeout)
	if err != nil {
		return nil, err
	}
	numPayloads := membuffers.GetUint32(sizeBuffer)

	t.logger.Info("receiving transport data", log.Int("payloads", int(numPayloads)), log.String("peer", conn.RemoteAddr().String()), trace.LogFieldFrom(ctx))

	if numPayloads > MAX_PAYLOADS_IN_MESSAGE {
		return nil, errors.Errorf("received message with too many payloads: %d", numPayloads)
	}

	for i := uint32(0); i < numPayloads; i++ {
		// receive payload size
		sizeBuffer, err := readTotal(ctx, conn, 4, timeout)
		if err != nil {
			return nil, err
		}
		payloadSize := membuffers.GetUint32(sizeBuffer)
		if payloadSize > MAX_PAYLOAD_SIZE_BYTES {
			return nil, errors.Errorf("received message with a payload too big: %d bytes", payloadSize)
		}

		// receive payload data
		payload, err := readTotal(ctx, conn, payloadSize, timeout)
		if err != nil {
			return nil, err
		}
		res = append(res, payload)

		// receive padding
		paddingSize := calcPaddingSize(uint32(len(payload)))
		if paddingSize > 0 {
			_, err := readTotal(ctx, conn, paddingSize, timeout)
			if err != nil {
				return nil, err
			}
		}
	}

	return res, nil
}

func (t *directTransport) notifyListener(ctx context.Context, payloads [][]byte) {
	listener := t.getListener()

	if listener == nil {
		return
	}

	listener.OnTransportMessageReceived(ctx, payloads)
}

func (t *directTransport) getListener() TransportListener {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.transportListenerUnderMutex
}

func (t *directTransport) clientMainLoop(parentCtx context.Context, address string, msgs chan *TransportData) {
	for {
		ctx := trace.NewContext(parentCtx, fmt.Sprintf("Gossip.Transport.TCP.Client.%s", address))
		t.logger.Info("attempting outgoing transport connection", log.String("server", address), trace.LogFieldFrom(ctx))
		conn, err := net.Dial("tcp", address)

		if err != nil {
			t.logger.Info("cannot connect to gossip peer endpoint", log.String("peer", address), trace.LogFieldFrom(ctx))
			time.Sleep(t.config.GossipConnectionKeepAliveInterval())
			continue
		}

		if !t.clientHandleOutgoingConnection(ctx, conn, msgs) {
			return
		}
	}
}

// returns true if should attempt reconnect on error
func (t *directTransport) clientHandleOutgoingConnection(ctx context.Context, conn net.Conn, msgs chan *TransportData) bool {
	t.logger.Info("successful outgoing gossip transport connection", log.String("peer", conn.RemoteAddr().String()), trace.LogFieldFrom(ctx))

	for {
		select {
		case data := <-msgs:
			err := t.sendTransportData(ctx, conn, data)
			if err != nil {
				t.metrics.outgoingConnectionFailedSend.Inc()
				t.logger.Info("failed sending transport data, reconnecting", log.Error(err), log.String("peer", conn.RemoteAddr().String()), trace.LogFieldFrom(ctx))
				conn.Close()
				return true
			}
		case <-time.After(t.config.GossipConnectionKeepAliveInterval()):
			err := t.sendKeepAlive(ctx, conn)
			if err != nil {
				t.metrics.outgoingConnectionFailedKeepalive.Inc()
				t.logger.Info("failed sending keepalive, reconnecting", log.Error(err), log.String("peer", conn.RemoteAddr().String()), trace.LogFieldFrom(ctx))
				conn.Close()
				return true
			}
		case <-ctx.Done():
			t.logger.Info("client loop stopped since server is shutting down", trace.LogFieldFrom(ctx))
			conn.Close()
			return false
		}
	}
}

func (t *directTransport) sendTransportData(ctx context.Context, conn net.Conn, data *TransportData) error {
	t.logger.Info("sending transport data", log.Int("payloads", len(data.Payloads)), log.String("peer", conn.RemoteAddr().String()), trace.LogFieldFrom(ctx))

	timeout := t.config.GossipNetworkTimeout()
	zeroBuffer := make([]byte, 4)
	sizeBuffer := make([]byte, 4)

	// send num payloads
	membuffers.WriteUint32(sizeBuffer, uint32(len(data.Payloads)))
	err := write(ctx, conn, sizeBuffer, timeout)
	if err != nil {
		return err
	}

	for _, payload := range data.Payloads {
		// send payload size
		membuffers.WriteUint32(sizeBuffer, uint32(len(payload)))
		err := write(ctx, conn, sizeBuffer, timeout)
		if err != nil {
			return err
		}

		// send payload data
		err = write(ctx, conn, payload, timeout)
		if err != nil {
			return err
		}

		// send padding
		paddingSize := calcPaddingSize(uint32(len(payload)))
		if paddingSize > 0 {
			err = write(ctx, conn, zeroBuffer[:paddingSize], timeout)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func calcPaddingSize(size uint32) uint32 {
	const contentAlignment = 4
	alignedSize := (size + contentAlignment - 1) / contentAlignment * contentAlignment
	return alignedSize - size
}

func (t *directTransport) sendKeepAlive(ctx context.Context, conn net.Conn) error {
	t.logger.Info("sending transport data", log.Int("payloads", 0), log.String("peer", conn.RemoteAddr().String()), trace.LogFieldFrom(ctx))

	timeout := t.config.GossipNetworkTimeout()
	zeroBuffer := make([]byte, 4)

	// send zero num payloads
	err := write(ctx, conn, zeroBuffer, timeout)
	if err != nil {
		return err
	}

	return nil
}

func readTotal(ctx context.Context, conn net.Conn, totalSize uint32, timeout time.Duration) ([]byte, error) {
	// TODO(https://github.com/orbs-network/orbs-network-go/issues/182): consider whether the right approach is to poll context this way or have a single watchdog goroutine that closes all active connections when context is cancelled
	// make sure context is still open
	err := ctx.Err()
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, totalSize)
	totalRead := uint32(0)
	for totalRead < totalSize {
		err := conn.SetReadDeadline(time.Now().Add(timeout))
		if err != nil {
			return nil, err
		}
		read, err := conn.Read(buffer[totalRead:])
		totalRead += uint32(read)
		if totalRead == totalSize {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return buffer, nil
}

func write(ctx context.Context, conn net.Conn, buffer []byte, timeout time.Duration) error {
	// TODO(https://github.com/orbs-network/orbs-network-go/issues/182): consider whether the right approach is to poll context this way or have a single watchdog goroutine that closes all active connections when context is cancelled
	// make sure context is still open
	err := ctx.Err()
	if err != nil {
		return err
	}

	err = conn.SetWriteDeadline(time.Now().Add(timeout))
	if err != nil {
		return err
	}
	written, err := conn.Write(buffer)
	if written != len(buffer) {
		if err == nil {
			return errors.Errorf("attempted to write %d bytes but only wrote %d", len(buffer), written)
		} else {
			return errors.Wrapf(err, "attempted to write %d bytes but only wrote %d", len(buffer), written)
		}
	}
	return nil
}
