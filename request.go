package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ffiat/nostr"
)

func NewRequest(cfg *Config) *Request {

	gc := &Request{
		fs:  flag.NewFlagSet("req", flag.ContinueOnError),
		cfg: cfg,
	}

	gc.fs.StringVar(&gc.npub, "npub", "", "event text note of Kind 1")
	gc.fs.BoolVar(&gc.following, "following", false, "event text note of Kind 1")

	return gc
}

type Request struct {
	fs  *flag.FlagSet
	cfg *Config

	npub      string
	following bool
}

func (g *Request) Name() string {
	return g.fs.Name()
}

func (g *Request) Init(args []string) error {
	return g.fs.Parse(args)
}

func (s *Request) Run() error {

	if s.npub != "" {

		pk, err := nostr.DecodeBech32(s.npub)
		if err != nil {
			log.Fatalf("\nunable to decode npub: %#v", err)
		}

		// List only the latest 3 event from the author.
		f := nostr.Filter{
			Authors: []string{pk.(string)},
			Kinds:   []uint32{nostr.KindTextNote},
			Limit:   10,
		}

		for _, v := range s.cfg.Relays {

			//log.Printf("Requesting from relay: %s", k)

			cc := NewConnection(v)
			err := cc.Listen()
			if err != nil {
				log.Fatalf("unable to listen to relay: %v", err)
			}
			defer cc.Close()

			sub, err := cc.Subscribe(nostr.Filters{f})
			if err != nil {
				log.Fatalf("\nunable to subscribe: %#v", err)
			}

			// FIXME: This is probabily a race condition

			time.Sleep(1 * time.Second)

			for event := range sub.EventStream {
				fmt.Printf("  [%s]\n\n", event.CreatedAt.Time())
				fmt.Printf("    ⤷  %s\n\n", event.Content)
			}
		}
	}

	if s.following {

		// 1. Range over following as defined in local config.

		config, err := DecodeConfig(s.cfg.Path)
		if err != nil {
			log.Fatalf("unable to decode local config: %v", err)
		}

		for _, v := range s.cfg.Relays {

			//log.Printf("Requesting from relay: %s", k)

			cc := NewConnection(v)
			err := cc.Listen()
			if err != nil {
				log.Fatalf("unable to listen to relay: %v", err)
			}

			for _, author := range config.Following {

				pk, err := nostr.DecodeBech32(author.PublicKey)
				if err != nil {
					log.Fatalf("\nunable to decode npub: %#v", err)
				}

				// List only the latest 3 event from the author.
				f := nostr.Filter{
					Authors: []string{pk.(string)},
					Kinds:   []uint32{nostr.KindTextNote},
					Limit:   10,
				}

				sub, err := cc.Subscribe(nostr.Filters{f})
				if err != nil {
					log.Fatalf("\nunable to subscribe: %#v", err)
				}

				// FIXME: This is probabily a race condition.
				// I think I fixed this with the ordone pattern.
				//time.Sleep(2 * time.Second)
				orDone := func(done <-chan struct{}, c <-chan *nostr.Event) <-chan *nostr.Event {
					valStream := make(chan *nostr.Event)
					go func() {
						defer close(valStream)
						for {
							select {
							case <-done:
								return
							case v, ok := <-c:
								if ok == false {
									return
								}
								valStream <- v
							}
						}
					}()
					return valStream
				}

				for e := range orDone(sub.Done, sub.EventStream) {
					ShowEvent(e, author)
				}
			}

			cc.Close()
		}
	}

	return nil
}
