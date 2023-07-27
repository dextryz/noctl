package main

import (
	"flag"
	"log"
)

func NewFollow(cfg *Config) *Follow {

	gc := &Follow{
		fs:  flag.NewFlagSet("follow", flag.ContinueOnError),
		cfg: cfg,
	}

	gc.fs.StringVar(&gc.ls, "ls", "", "list all following users")
	gc.fs.StringVar(&gc.add, "add", "", "user public key to add to following list")
	gc.fs.StringVar(&gc.remove, "remove", "", "remove user via public key")

	return gc
}

type Follow struct {
	fs  *flag.FlagSet
	cfg *Config

	ls     string
	add    string
	remove string
}

func (g *Follow) Name() string {
	return g.fs.Name()
}

func (g *Follow) Init(args []string) error {
	return g.fs.Parse(args)
}

func (s *Follow) Run() error {

	if s.add != "" {
		s.cfg.Following[s.add] = Author{
			PublicKey: s.add,
		}
		s.cfg.Save()
	}

	if s.remove != "" {
		log.Println("[remove] not implemented")
	}

	if s.ls != "" {
		log.Println("[ls] not implemented")
	}

	return nil
}
