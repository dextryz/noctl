package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/ffiat/nostr"
)

func NewRelay(cfg *Config) *Relay {

	gc := &Relay{
		fs:  flag.NewFlagSet("relay", flag.ContinueOnError),
		cfg: cfg,
	}

	gc.fs.StringVar(&gc.info, "info", "", "event text note of Kind 1")
	gc.fs.StringVar(&gc.add, "add", "", "event text note of Kind 1")
	gc.fs.StringVar(&gc.remove, "remove", "", "event text note of Kind 2")

	return gc
}

type Relay struct {
	fs  *flag.FlagSet
	cfg *Config

	info   string
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

	if s.info != "" {

		ctx := context.Background()

		// normalize URL to start with http:// or https://
		if !strings.HasPrefix(s.info, "http") && !strings.HasPrefix(s.info, "ws") {
			s.info = "wss://" + s.info
		}
		p, err := url.Parse(s.info)
		if err != nil {
			return fmt.Errorf("Cannot parse url: %s", s.info)
		}
		if p.Scheme == "ws" {
			p.Scheme = "http"
		} else if p.Scheme == "wss" {
			p.Scheme = "https"
		}
		p.Path = strings.TrimRight(p.Path, "/")

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.String(), nil)

		// add the NIP-11 header
		req.Header.Add("Accept", "application/nostr+json")

		// send the request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		info := &nostr.RelayInformation{}
		err = json.NewDecoder(resp.Body).Decode(info)
		if err != nil {
			return err
		}

		PrintJson(info)
	}

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
