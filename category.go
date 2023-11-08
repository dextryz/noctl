package main

import (
	"flag"
	"log"
	"strings"

	"github.com/ffiat/nostr"
)

func NewCategory(cfg *Config) *Category {

	gc := &Category{
		fs:  flag.NewFlagSet("category", flag.ContinueOnError),
		cfg: cfg,
	}

	gc.fs.StringVar(&gc.name, "name", "", "event text note of Kind 1")
	gc.fs.StringVar(&gc.pubkeys, "pubkeys", "", "event text note of Kind 1")

	return gc
}

type Category struct {
	fs  *flag.FlagSet
	cfg *Config

	name    string
	pubkeys string
}

func (g *Category) Name() string {
	return g.fs.Name()
}

func (g *Category) Init(args []string) error {
	return g.fs.Parse(args)
}

func (s *Category) Run() error {

	if s.name == "" {
		log.Fatalln("please give this category a name")
	}

	if s.pubkeys == "" {
		log.Fatalln("please give at least one pubkey")
	}

	pubkeys := strings.Split(s.pubkeys, ",")

	tags := nostr.Tags{}

	tags = append(tags, nostr.Tag{"d", s.name})

	for _, key := range pubkeys {

		_, pk, err := nostr.DecodeBech32(key)
		if err != nil {
			log.Fatalln(err)
		}

		t := nostr.Tag{"p", pk}
		tags = append(tags, t)
	}

	log.Println("TAGS")
	log.Println(tags)

	e := nostr.Event{
		// TODO: Update to nostr.CategoryPeople
		Kind:      3000,
		Tags:      tags,
		CreatedAt: nostr.Now(),
		Content:   "",
	}

	for k, v := range s.cfg.Relays {

		log.Printf("Publishing category to relay: %s", k)

		cc := NewConnection(v)
		err := cc.Listen()
		if err != nil {
			log.Fatalf("unable to listen to relay: %v", err)
		}
		defer cc.Close()

		msg, err := cc.Publish(e, s.cfg.PrivateKey)
		if msg.Ok {
			log.Println("[\033[1;32m+\033[0m] Category note published")
			log.Println(msg)
		} else {
			// TODO:
			log.Println("What")
			log.Println(msg)
			log.Println(err)
		}
		if err != nil {
			return err
		}
	}

	return nil
}
