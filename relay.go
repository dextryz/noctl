package main

import (
	"flag"
	"log"
)

func NewRelay(cfg *Config) *Relay {

	gc := &Relay{
		fs:  flag.NewFlagSet("relay", flag.ContinueOnError),
		cfg: cfg,
	}

	gc.fs.StringVar(&gc.add, "add", "", "event text note of Kind 1")
	gc.fs.StringVar(&gc.remove, "remove", "", "event text note of Kind 2")

	return gc
}

type Relay struct {
	fs  *flag.FlagSet
	cfg *Config

	add    string
	remove string
}

func (g *Relay) Name() string {
	return g.fs.Name()
}

func (g *Relay) Init(args []string) error {
	return g.fs.Parse(args)
}

func (s *Relay) Run() error {

	if s.add != "" {
		s.cfg.AddRelay(s.add)
		s.cfg.Save()
	}

	if s.remove != "" {
		s.cfg.RemoveRelay(s.add)
		s.cfg.Save()
	}

	return nil
}

// View the current state of the profile as defined in CONFIG_PATH
func (s *Relay) view() error {

	config, err := DecodeConfig(s.cfg.Path)
	if err != nil {
		log.Fatalf("unable to decode local config: %v", err)
	}

	PrintJson(config)

	return nil
}
