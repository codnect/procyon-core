package runtime

import "context"

type CommandLineRunner interface {
	Run(ctx context.Context, args *Arguments) error
}
