package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"time"

	"git.net.quant.sh/tools/geth-healthcheck/internal/checker"
	"git.net.quant.sh/tools/geth-healthcheck/internal/nodegetter"
	"github.com/alecthomas/kong"
)

var CLI struct {
	NodeURL         *url.URL `kong:"required,help='The URL of the eth not to be monitored',env='NODE_URL',default='http://127.0.0.1:8545'"`
	ExternalNodeURL *url.URL `kong:"required,help='The URL of the eth for block comparison',env='EXTERNAL_NODE_URL'"`

	MaxBlockDifference uint          `kong:"env='MAX_BLOCK_DIFFERENCE',default='3'"`
	MaxNodeBlockAge    time.Duration `kong:"env='MAX_BLOCK_AGE',default='30s'"`

	Listen string `kong:"env='LISTEN',default='127.0.0.1:8080'"`
}

func main() {
	kong.Parse(&CLI)

	ctx := context.Background()

	nodeGetter, err := nodegetter.NewNodeGetter(ctx, CLI.NodeURL.String())
	if err != nil {
		log.Printf("ERROR: could not connect to node: %v", err)
	}

	externalNodeGetter, err := nodegetter.NewNodeGetter(ctx, CLI.ExternalNodeURL.String())
	if err != nil {
		log.Printf("ERROR: could not connect to external node: %v", err)
	}

	checkerConf := checker.Config{
		MaxBlockDifference:  CLI.MaxBlockDifference,
		MaxNodeBlockAge:     CLI.MaxNodeBlockAge,
		NodeBlockGetter:     nodeGetter,
		ExternalBlockGetter: externalNodeGetter,
	}

	checker := checker.NewChecker(checkerConf)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("checking node")
		valid, err := checker.IsHealthy(ctx)
		if err != nil {
			log.Printf("ERROR: could not verrify node; %v", err)
		}

		if valid {
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusBadGateway)
	})

	log.Fatal(http.ListenAndServe(CLI.Listen, nil))
}
