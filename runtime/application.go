package runtime

type Application interface {
	Run(args ...string) error
	Exit() int
}
