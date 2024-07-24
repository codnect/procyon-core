package main

import (
	"context"
	"github.com/codnect/procyoncore/component"
	"github.com/codnect/procyoncore/component/filter"
)

func main() {

	component.Register(NewUserController, component.Named("test"))
	component.Register(NewUserService)

	objectContainer := component.NewObjectContainer()
	for _, cmp := range component.List() {
		objectContainer.Definitions().Register(cmp.Definition())
	}

	instance, err := objectContainer.GetObject(context.Background(), filter.ByTypeOf[*UserController](), filter.ByName("23"))

	if err != nil {
		panic(err)
	}

	if instance != nil {

	}
}

type UserController struct {
	userService *UserService
}

func (c *UserController) saveUser() {

}

func NewUserController(userService *UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}
