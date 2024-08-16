package runtime

type Application interface {
	Context() Context
	Run(args ...string) error
}
