package main

import (
	"context"
	"github.com/codnect/procyoncore/component"
	"github.com/codnect/procyoncore/component/filter"
	"github.com/codnect/procyoncore/runtime"
	"time"
)

func main() {

	component.Register(NewUserController, component.Named("test"))
	component.Register(NewUserController, component.Named("test2"))
	component.Register(NewUserService)

	objectContainer := component.NewObjectContainer()
	for _, cmp := range component.List() {
		objectContainer.Definitions().Register(cmp.Definition())
	}

	var app runtime.Application
	if app != nil {

	}
	lst := objectContainer.ListObjects(context.Background(), filter.ByTypeOf[*UserController](), filter.ByName("test"))
	if len(lst) != 0 {

	}

	x, err := objectContainer.GetObject(context.Background(), filter.ByName("procyonRuntimeCustomizer"))
	if x != nil {

	}
	if err != nil {

	}
}

type MyEvent struct {
}

func (c MyEvent) EventSource() any {
	return nil
}

func (c MyEvent) Time() time.Time {
	return time.Time{}
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
