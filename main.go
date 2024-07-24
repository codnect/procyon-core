package main

import (
	"github.com/codnect/procyoncore/component"
	"github.com/codnect/procyoncore/component/filter"
)

func main() {

	component.Register(NewController, component.Named("test"), component.Scoped("myscope")).
		ConditionalOn(nil)

	container := component.NewObjectContainer()
	listOfComponents := component.List()
	// register definitions ....
	for _, cm := range listOfComponents {
		container.Singletons().
	}

	c, err := container.GetObject(nil, filter.ByName("test"), filter.ByType[string]())

}

type HelloService struct {
}

type HelloController struct {
}

func NewController(helloService *HelloService) *HelloController {
	return &HelloController{}
}

type UserController struct {
}

func NewUserController(helloService *HelloService) *HelloController {
	return &HelloController{}
}
