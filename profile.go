package main

import (
	"flag"
	"log"

	"github.com/dextryz/nostr"
)

func NewProfile(cfg *Config) *Profile {

	gc := &Profile{
		fs:  flag.NewFlagSet("profile", flag.ContinueOnError),
		cfg: cfg,
	}

	gc.fs.StringVar(&gc.nsec, "nsec", "", "event text note of Kind 1")
	gc.fs.StringVar(&gc.name, "name", "", "event text note of Kind 1")
	gc.fs.StringVar(&gc.about, "about", "", "event text note of Kind 0")
	gc.fs.StringVar(&gc.picture, "picture", "", "event text note of Kind 2")

	gc.fs.BoolVar(&gc.show, "show", false, "event text note of Kind 1")
	gc.fs.BoolVar(&gc.commit, "commit", false, "event text note of Kind 1")

	return gc
}

type Profile struct {
	fs  *flag.FlagSet
	cfg *Config

	nsec    string
	name    string
	about   string
	picture string
	show    bool
	commit  bool
}

func (g *Profile) Name() string {
	return g.fs.Name()
}

func (g *Profile) Init(args []string) error {
	return g.fs.Parse(args)
}

func (s *Profile) Run() error {

	if s.show {
		s.view()
	}

	// Decode the NIP-19 secret to a NIP-01 secret.
	if s.nsec != "" {

		var sk string
		if _, value, e := nostr.DecodeBech32(s.nsec); e == nil {
			sk = value
		}

		// Update secret
		s.cfg.PrivateKey = sk

		// Set public with which the event wat pushed.
		if pub, e := nostr.GetPublicKey(sk); e == nil {
			s.cfg.PublicKey = pub
		}

		s.cfg.Save()
	}

	if s.name != "" {
		s.cfg.Profile.Name = s.name
		s.cfg.Save()
	}

	if s.about != "" {
		s.cfg.Profile.About = s.about
		s.cfg.Save()
	}

	if s.picture != "" {
		s.cfg.Profile.Picture = s.picture
		s.cfg.Save()
	}

	// Commit event to relays to update profile.
	if s.commit {

		if len(s.cfg.Relays) == 0 {
			log.Fatalln("please add a relay before commiting")
		}

		for k, v := range s.cfg.Relays {

			log.Printf("Commiting profile changes to relay: %s", k)

			cc := NewConnection(v)
			err := cc.Listen()
			if err != nil {
				log.Fatalf("unable to listen to relay: %v", err)
			}
			defer cc.Close()

			e := nostr.Event{
				Kind:      nostr.KindSetMetadata,
				Tags:      nil,
				CreatedAt: nostr.Now(),
				Content:   s.cfg.Profile.String(),
			}

			_, err = cc.Publish(e, s.cfg.PrivateKey)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// View the current state of the profile as defined in CONFIG_PATH
func (s *Profile) view() error {

	config, err := DecodeConfig(s.cfg.Path)
	if err != nil {
		log.Fatalf("unable to decode local config: %v", err)
	}

	PrintJson(config)

	return nil
}
