package checker

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type fakeBlockGetter struct {
	block    Block
	hasError bool
}

func (r *fakeBlockGetter) GetCurrentBlock(ctx context.Context) (Block, error) {
	if r.hasError {
		return Block{}, fmt.Errorf("has error")
	}

	return r.block, nil

}

func TestChecker_IsHealthy(t *testing.T) {
	now := time.Now()

	type fields struct {
		maxBlockDifference  uint
		maxNodeBlockAge     time.Duration
		externalBlock       Block
		nodeBlock           Block
		externalBlockGetter BlockGetter
		nodeBlockGetter     BlockGetter
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "healthy",
			fields: fields{
				maxBlockDifference: 0,
				externalBlockGetter: &fakeBlockGetter{
					block: Block{Number: 1, Timestamp: now},
				},
				nodeBlockGetter: &fakeBlockGetter{
					block: Block{Number: 1, Timestamp: now},
				},
			},
			args:    args{ctx: context.Background()},
			want:    true,
			wantErr: false,
		},
		{
			name: "external block higher",
			fields: fields{
				maxBlockDifference: 0,
				externalBlockGetter: &fakeBlockGetter{
					block: Block{Number: 2, Timestamp: now.Add(time.Minute)},
				},
				nodeBlockGetter: &fakeBlockGetter{
					block: Block{Number: 1, Timestamp: now},
				},
			},
			args:    args{ctx: context.Background()},
			want:    false,
			wantErr: false,
		},
		{
			name: "internal block higher",
			fields: fields{
				maxBlockDifference: 0,
				externalBlockGetter: &fakeBlockGetter{
					block: Block{Number: 1, Timestamp: time.Unix(1618260129, 0)},
				},
				nodeBlockGetter: &fakeBlockGetter{
					block: Block{Number: 2, Timestamp: time.Unix(1618260129, 0)},
				},
			},
			args:    args{ctx: context.Background()},
			want:    true,
			wantErr: false,
		},
		{
			name: "block within range",
			fields: fields{
				maxBlockDifference: 2,
				externalBlockGetter: &fakeBlockGetter{
					block: Block{Number: 2, Timestamp: time.Unix(1618260129, 0)},
				},
				nodeBlockGetter: &fakeBlockGetter{
					block: Block{Number: 0, Timestamp: time.Unix(1618260129, 0)},
				},
			},
			args:    args{ctx: context.Background()},
			want:    true,
			wantErr: false,
		},
		{
			name: "external unavailable internal block not stale",
			fields: fields{
				maxBlockDifference: 2,
				maxNodeBlockAge:    time.Second * 10,
				externalBlockGetter: &fakeBlockGetter{
					block:    Block{},
					hasError: true,
				},
				nodeBlockGetter: &fakeBlockGetter{
					block: Block{Number: 0, Timestamp: now},
				},
			},
			args:    args{ctx: context.Background()},
			want:    true,
			wantErr: false,
		},
		{
			name: "external unavailable internal block stale",
			fields: fields{
				maxBlockDifference: 2,
				maxNodeBlockAge:    time.Second * 10,
				externalBlockGetter: &fakeBlockGetter{
					block:    Block{},
					hasError: true,
				},
				nodeBlockGetter: &fakeBlockGetter{
					block: Block{Number: 5, Timestamp: now.Add(-time.Minute * 2)},
				},
			},
			args:    args{ctx: context.Background()},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Checker{
				maxBlockDifference:  tt.fields.maxBlockDifference,
				maxNodeBlockAge:     tt.fields.maxNodeBlockAge,
				externalBlock:       tt.fields.externalBlock,
				nodeBlock:           tt.fields.nodeBlock,
				externalBlockGetter: tt.fields.externalBlockGetter,
				nodeBlockGetter:     tt.fields.nodeBlockGetter,
			}
			got, err := c.IsHealthy(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Checker.IsHealthy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Checker.IsHealthy() = %v, want %v", got, tt.want)
			}
		})
	}
}
