package retry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestRetryHandle(t *testing.T) {
	type args struct {
		ctx      context.Context
		base     time.Duration
		maxRetry uint64
		f        RetryFunc
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				ctx:      context.Background(),
				base:     1 * time.Second,
				maxRetry: 2,
				f: func(ctx context.Context) error {
					fmt.Println("test1")
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				ctx:      context.Background(),
				base:     1 * time.Second,
				maxRetry: 2,
				f: func(ctx context.Context) error {
					fmt.Println("test2")
					return RetryableError(errors.New("error"))
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RetryHandle(tt.args.ctx, tt.args.base, tt.args.maxRetry, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("RetryHandle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
