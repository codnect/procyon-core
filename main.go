package main

import "github.com/codnect/procyoncore/component"

func main() {

	component.Register(NewHelloController).Named("homeController").
		Scoped(component.SingletonScope).
		Prioritized(2).
		Inject(0).
		ConditionalOn(nil)

	l := component.List()
	if len(l) != 0 {

	}
	d, _ := component.MakeDefinition(NewHelloController,
		component.WithName(""), component.WithNamedArgument(0, ""))
	if d.IsPrimary() {

	}

}

type HelloController struct {
}

func NewHelloController(name string) *HelloController {
	return &HelloController{}
}
