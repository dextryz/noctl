package main

import (
	"flag"
	"log"
	"os"
)

var (
	CONFIG_NOSTR = StringEnv("CONFIG_NOSTR")
)

type Runner interface {
	Init([]string) error
	Run() error
	Name() string
}

func main() {

	flag.Parse()
	log.SetFlags(0)

	cfg, err := DecodeConfig(CONFIG_NOSTR)
	if err != nil {
		log.Fatalf("unable to decode local cfg: %v", err)
	}

	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatalln(err)
	}

	cmds := []Runner{
		NewRelay(cfg),
		NewProfile(cfg),
		NewEvent(cfg),
		NewRequest(cfg),
		NewFollow(cfg),
		NewCategory(cfg),
		NewEncode(cfg),
		NewKey(cfg),
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			cmd.Init(os.Args[2:])
			err := cmd.Run()
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}
