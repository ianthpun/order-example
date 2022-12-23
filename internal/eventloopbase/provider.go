package eventloopbase

import "context"

type Provider interface {
	Start(ctx context.Context) error
}
