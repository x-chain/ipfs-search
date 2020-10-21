package eventsource

import (
	"context"
	"fmt"
	"log"

	"github.com/ipfs-search/ipfs-sniffer/proxy"

	"github.com/ipfs/go-datastore"
	"github.com/libp2p/go-eventbus"
	"github.com/libp2p/go-libp2p-core/event"
)

const bufSize = 256

type EventSource struct {
	bus     event.Bus
	emitter event.Emitter
	ds      datastore.Batching
}

func New(b event.Bus, ds datastore.Batching) (EventSource, error) {
	e, err := b.Emitter(new(EvtProviderPut))
	if err != nil {
		return EventSource{}, err
	}

	s := EventSource{
		bus:     b,
		emitter: e,
	}

	s.ds = proxy.New(ds, s.afterPut)

	return s, nil
}

// nonFatalError is called on non-fatal errors
func (s *EventSource) nonFatalError(err error) {
	log.Printf("error: %v\n", err)
}

func (s *EventSource) afterPut(k datastore.Key, v []byte, err error) error {
	// Ignore error'ed Put's
	if err != nil {
		return err
	}

	// Ignore non-provider keys
	if !isProviderKey(k) {
		return nil
	}

	cid, err := keyToCID(k)
	if err != nil {
		s.nonFatalError(fmt.Errorf("cid from key '%s': %w", k, err))
		return nil
	}

	pid, err := keyToPeerID(k)
	if err != nil {
		s.nonFatalError(fmt.Errorf("pid from key '%s': %w", k, err))
		return nil
	}

	err = s.emitter.Emit(EvtProviderPut{
		CID:    cid,
		PeerID: pid,
	})
	if err != nil {
		s.nonFatalError(fmt.Errorf("cid from key '%s': %w", k, err))
		return nil
	}

	return nil
}

func (s *EventSource) Batching() datastore.Batching {
	return s.ds
}

// Subscribe handleFunc to EvtProviderPut events
func (s *EventSource) Subscribe(ctx context.Context, handleFunc func(context.Context, EvtProviderPut) error) error {
	sub, err := s.bus.Subscribe(new(EvtProviderPut), eventbus.BufSize(bufSize))
	if err != nil {
		return fmt.Errorf("subscribing: %w", err)
	}
	defer sub.Close()

	c := sub.Out()
	for {
		select {
		case <-ctx.Done():
			return err
		case e, ok := <-c:
			if !ok {
				return fmt.Errorf("reading from event bus")
			}

			evt, ok := e.(EvtProviderPut)
			if !ok {
				return fmt.Errorf("casting event: %v", evt)
			}

			err := handleFunc(ctx, evt)
			if err != nil {
				return err
			}
		}
	}
}