package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/dextryz/nostr"
)

func NewKey(cfg *Config) *Key {

	gc := &Key{
		fs:  flag.NewFlagSet("key", flag.ContinueOnError),
		cfg: cfg,
	}

	gc.fs.StringVar(&gc.key, "decode", "", "event text note of Kind 1")
	gc.fs.BoolVar(&gc.gen, "new", false, "event text note of Kind 1")

	return gc
}

type Key struct {
	fs  *flag.FlagSet
	cfg *Config

	key string
	gen bool
}

func (g *Key) Name() string {
	return g.fs.Name()
}

func (g *Key) Init(args []string) error {
	return g.fs.Parse(args)
}

func (s *Key) Run() error {

	// TODO: For now just use the raw event. Later implement NIP-21
	if s.key != "" {

		_, pubkey, err := nostr.DecodeBech32(s.key)
		if err != nil {
			log.Fatal("unable to generate public key")
		}
		fmt.Printf("%s", pubkey)
	}

	if s.gen {

		sk := nostr.GeneratePrivateKey()
		pk, err := nostr.GetPublicKey(sk)

		if err != nil {
			log.Fatal("unable to generate public key")
		}

		ns, err := nostr.EncodePrivateKey(sk)
		if err != nil {
			log.Fatal("unable to generate public key")
		}

		np, err := nostr.EncodePublicKey(pk)
		if err != nil {
			log.Fatal("unable to generate public key")
		}

		fmt.Printf("nsec: %s\n", ns)
		fmt.Printf("npub: %s", np)
	}

	return nil
}
