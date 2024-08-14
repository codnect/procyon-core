package runtime

type Application interface {
	Run(args ...string) error
	Exit() int
}

type ApplicationRunner interface {
	Run(ctx Context, args *Arguments) error
}
