package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/dextryz/nostr"
)

// > nix req -npub npub15pa38uvf6vyux63en5sgyzxsnp3myvwrw3dan43afhjr8x44gzjse75plm -kind 5001

func NewRequest(cfg *Config) *Request {

	gc := &Request{
		fs:  flag.NewFlagSet("req", flag.ContinueOnError),
		cfg: cfg,
	}

	gc.fs.StringVar(&gc.id, "id", "", "event ID of text note")
	gc.fs.StringVar(&gc.npub, "npub", "", "event text note of Kind 1")
	gc.fs.IntVar(&gc.kind, "kind", 1, "set event kind to be pulled from relays")
	gc.fs.BoolVar(&gc.following, "following", false, "event text note of Kind 1")

	return gc
}

type Request struct {
	fs  *flag.FlagSet
	cfg *Config

	id        string
	npub      string
	kind      int
	following bool
}

func (g *Request) Name() string {
	return g.fs.Name()
}

func (g *Request) Init(args []string) error {
	return g.fs.Parse(args)
}

func (s *Request) Run() error {

	// TODO: For now just use the raw event. Later implement NIP-21
	if s.id != "" {

		f := nostr.Filter{
			Ids:   []string{s.id},
			Kinds: []uint32{nostr.KindTextNote},
			Limit: 10,
		}

		for _, v := range s.cfg.Relays {

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

			// This should be a single event
			if len(sub.EventStream) != 1 {
				//				log.Fatalf("more than one event was pulled: %d", len(sub.EventStream))
			}

			for event := range sub.EventStream {
				fmt.Printf("  [%s]\n\n", event.CreatedAt.Time())
				fmt.Printf("    └──  %s\n\n", event.Content)
			}
		}
	}

	if s.npub != "" {

		_, pk, err := nostr.DecodeBech32(s.npub)
		if err != nil {
			log.Fatalf("\nunable to decode npub: %#v", err)
		}

		// List only the latest 3 event from the author.
		f := nostr.Filter{
			Authors: []string{pk},
			//Kinds:   []uint32{nostr.KindTextNote},
			Kinds: []uint32{uint32(s.kind)},
			Limit: 10,
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

			log.Printf("%d events found for %s", len(sub.EventStream), s.npub)

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
				PrintJson(e)
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

				_, pk, err := nostr.DecodeBech32(author.PublicKey)
				if err != nil {
					log.Fatalf("\nunable to decode npub: %#v", err)
				}

				// List only the latest 3 event from the author.
				f := nostr.Filter{
					Authors: []string{pk},
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
