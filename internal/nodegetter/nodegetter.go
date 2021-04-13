package nodegetter

import (
	"context"
	"time"

	"git.net.quant.sh/tools/geth-healthcheck/internal/checker"
	"github.com/ethereum/go-ethereum/ethclient"
)

type NodeGetter struct {
	client *ethclient.Client
}

func NewNodeGetter(ctx context.Context, url string) (checker.BlockGetter, error) {
	client, err := ethclient.DialContext(ctx, url)
	if err != nil {
		return nil, err
	}

	return &NodeGetter{
		client: client,
	}, nil
}

func (g *NodeGetter) GetCurrentBlock(ctx context.Context) (checker.Block, error) {
	block, err := g.client.BlockByNumber(ctx, nil)
	if err != nil {
		return checker.Block{}, err
	}

	return checker.Block{
		Timestamp: time.Unix(int64(block.Time()), 0),
		Number:    int(block.NumberU64()),
	}, nil
}
