package checker

import (
	"context"
	"fmt"
	"time"
)

type Block struct {
	Timestamp time.Time
	Number    int
}

type BlockGetter interface {
	GetCurrentBlock(ctx context.Context) (Block, error)
}

type Config struct {
	MaxBlockDifference uint
	MaxNodeBlockAge    time.Duration

	ExternalBlockGetter BlockGetter
	NodeBlockGetter     BlockGetter
}

type Checker struct {
	maxBlockDifference uint
	maxNodeBlockAge    time.Duration

	externalBlock Block
	nodeBlock     Block

	externalBlockGetter BlockGetter
	nodeBlockGetter     BlockGetter
}

func NewChecker(c Config) *Checker {
	return &Checker{
		maxBlockDifference: c.MaxBlockDifference,
		maxNodeBlockAge:    c.MaxNodeBlockAge,

		nodeBlockGetter:     c.NodeBlockGetter,
		externalBlockGetter: c.ExternalBlockGetter,
	}
}

func (c *Checker) IsHealthy(ctx context.Context) (bool, error) {
	err := c.fetchNodeBlock(ctx)
	if err != nil {
		return false, err
	}

	blockValid := c.isNodeSynced(ctx)

	return blockValid, nil
}

func (c *Checker) fetchNodeBlock(ctx context.Context) error {
	currentNodeBlock, err := c.nodeBlockGetter.GetCurrentBlock(ctx)
	if err != nil {
		return fmt.Errorf("could not fetch latest node block: %v", err)
	}
	c.nodeBlock = currentNodeBlock
	return nil
}

func (c *Checker) isNodeSynced(ctx context.Context) bool {
	err := c.fetchExternalBlock(ctx)
	if err != nil {
		return !c.isNodeBlockStale()
	}

	return c.blocksAreInRange()
}

func (c *Checker) fetchExternalBlock(ctx context.Context) error {
	externalBlock, err := c.externalBlockGetter.GetCurrentBlock(ctx)
	if err != nil {
		fmt.Println("could not fetch latest external block:", err)
	}

	c.externalBlock = externalBlock
	return err
}

func (c *Checker) isNodeBlockStale() bool {
	return time.Since(c.nodeBlock.Timestamp) > c.maxNodeBlockAge
}

func (c *Checker) blocksAreInRange() bool {
	blockDiff := c.externalBlock.Number - c.nodeBlock.Number

	return blockDiff <= int(c.maxBlockDifference)
}
