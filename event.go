package main

import (
	"flag"
	"log"

	"github.com/ffiat/nostr"
)

func NewEvent(cfg *Config) *Event {

	gc := &Event{
		fs:  flag.NewFlagSet("event", flag.ContinueOnError),
		cfg: cfg,
	}

	gc.fs.StringVar(&gc.note, "note", "", "event text note of Kind 1")

	return gc
}

type Event struct {
	fs  *flag.FlagSet
	cfg *Config

	// Content of text note
	note string
	sign string
}

func (g *Event) Name() string {
	return g.fs.Name()
}

func (g *Event) Init(args []string) error {
	return g.fs.Parse(args)
}

func (s *Event) Run() error {

	if s.note != "" {

		e := nostr.Event{
			Kind:      nostr.KindTextNote,
			Tags:      nil,
			CreatedAt: nostr.Now(),
			Content:   s.note,
		}

		for k, v := range s.cfg.Relays {

			log.Printf("Publishing event to relay: %s", k)

			cc := NewConnection(v)
			err := cc.Listen()
			if err != nil {
				log.Fatalf("unable to listen to relay: %v", err)
			}
			defer cc.Close()

			ok, err := cc.Publish(e, s.cfg.PrivateKey)
			if ok != nil {
				log.Printf("[\033[1;32m+\033[0m] Text note published: [status: %v]", ok.Ok)
			}
			if err != nil {
				return err
			}
		}
	}

	return nil
}
