package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/dextryz/nostr"
)

func NewEncode(cfg *Config) *Encode {

	gc := &Encode{
		fs:  flag.NewFlagSet("encode", flag.ContinueOnError),
		cfg: cfg,
	}

	gc.fs.StringVar(&gc.nsec, "nsec", "", "encode a NIP-01 private key to a NIP-19 nsec")
	gc.fs.StringVar(&gc.npub, "npub", "", "encode a NIP-01 public key to a NIP-19 npub")
	gc.fs.StringVar(&gc.note, "note", "", "encode a NIP-01 event ID to a NIP-19 note")

	return gc
}

type Encode struct {
	fs  *flag.FlagSet
	cfg *Config

	nsec string
	npub string
	note string
}

func (g *Encode) Name() string {
	return g.fs.Name()
}

func (g *Encode) Init(args []string) error {
	return g.fs.Parse(args)
}

func (s *Encode) Run() error {

	if s.nsec != "" {
		key, err := nostr.EncodePrivateKey(s.nsec)
		if err != nil {
			log.Fatal("unable to encode secret key")
		}
		fmt.Println(key)
	} else if s.npub != "" {
		key, err := nostr.EncodePublicKey(s.npub)
		if err != nil {
			log.Fatal("unable to encode public key")
		}
		fmt.Println(key)
	} else if s.note != "" {
		key, err := nostr.EncodeNote(s.note)
		if err != nil {
			log.Fatalf("unable to encode note: %v", err)
		}
		fmt.Println(key)
	} else {
		log.Fatalln("unsupported key type")
	}

	return nil
}
