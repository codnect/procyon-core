package runtime

type Runner interface {
	Run(ctx Context, args *Arguments) error
}
