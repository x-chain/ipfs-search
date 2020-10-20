package commands

import (
	"context"
	"github.com/ipfs-search/ipfs-search/config"
	"github.com/ipfs-search/ipfs-search/crawler"
	"github.com/ipfs-search/ipfs-search/instr"

	"github.com/ipfs-search/ipfs-search/queue/amqp"
)

// AddHash queues a single IPFS hash for indexing
func AddHash(ctx context.Context, cfg *config.Config, hash string) error {
	instFlusher, err := instr.Install("ipfs-search-add")
	if err != nil {
		log.Fatal(err)
	}
	defer instFlusher()

	instrumentation := instr.New()

	f := amqp.PublisherFactory{
		AMQPURL: cfg.AMQP.AMQPURL,
		Queue:   "hashes",
	}

	queue, err := f.NewPublisher(ctx, instrumentation)
	if err != nil {
		return err
	}

	// Add with highest priority, as this is supposed to be available
	return queue.Publish(ctx, &crawler.Args{
		Hash: hash,
	}, 9)
}
