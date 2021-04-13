package nodegetter

import (
	"context"
	"testing"
)

func TestNodeGetter_GetCurrentBlock(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "infura",
			url:  "https://mainnet.infura.io/v3/b85fd340929e47c68f7f1531caf2667e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getter, err := NewNodeGetter(context.Background(), tt.url)
			if err != nil {
				t.Error(err)
			}

			block, err := getter.GetCurrentBlock(context.Background())
			if err != nil {
				t.Error(err)
			}

			t.Logf("Block number: %v, BlockTime: %v", block.Number, block.Timestamp)
		})
	}
}
