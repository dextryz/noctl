package main

import (
	"flag"
	"log"
	"strings"

	"github.com/dextryz/nostr"
)

// > nix summarize -note 'note1' -note 'note2' -length "1 phrase"

// StringSlice is a custom type for a slice of strings
type StringSlice []string

// String is required by the flag.Value interface
func (i *StringSlice) String() string {
	return strings.Join(*i, ",")
}

// Set is required by the flag.Value interface
func (i *StringSlice) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func NewSummarize(cfg *Config) *Summarize {

	gc := &Summarize{
		fs:  flag.NewFlagSet("summarize", flag.ContinueOnError),
		cfg: cfg,
	}

	gc.fs.Var(&gc.notes, "note", "List of notes (repeatable)")
	gc.fs.StringVar(&gc.length, "length", "", "event text note of Kind 1")

	return gc
}

type Summarize struct {
	fs  *flag.FlagSet
	cfg *Config

	notes  StringSlice
	length string
}

func (g *Summarize) Name() string {
	return g.fs.Name()
}

func (g *Summarize) Init(args []string) error {
	return g.fs.Parse(args)
}

func (s *Summarize) Run() error {

	if s.notes != nil {

		tags := nostr.Tags{
			nostr.Tag{"output", "text/plain"},
			nostr.Tag{"param", "length", s.length},
		}

		for _, v := range s.notes {
			tags = append(tags, nostr.Tag{"i", v, "text"})
		}

		e := nostr.Event{
			Kind:      nostr.KindSummarize,
			Content:   "",
			Tags:      tags,
			CreatedAt: nostr.Now(),
		}

		for k, v := range s.cfg.Relays {

			log.Printf("Publishing event to relay: %s", k)

			cc := NewConnection(v)
			err := cc.Listen()
			if err != nil {
				log.Fatalf("unable to listen to relay: %v", err)
			}
			defer cc.Close()

			msg, err := cc.Publish(e, s.cfg.PrivateKey)
			if msg.Ok {
				log.Println("[\033[1;32m+\033[0m] Text note published")
				//PrintJson(msg)
				log.Println(msg)
			} else {
				// TODO:
			}
			if err != nil {
				return err
			}
		}
	}

	return nil
}
