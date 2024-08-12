package main

import (
	"context"
	"github.com/codnect/procyoncore/component"
	"github.com/codnect/procyoncore/component/filter"
	"github.com/codnect/procyoncore/event"
	"github.com/codnect/procyoncore/runtime"
)

func main() {

	component.Register(NewUserController, component.Named("test"))
	component.Register(NewUserService)
	component.Register(newP)

	objectContainer := component.NewObjectContainer()

	for _, cmp := range component.List() {
		objectContainer.Definitions().Register(cmp.Definition())
	}

	err := objectContainer.Start(context.Background())

	if err != nil {
		panic(err)
	}

	var r runtime.Application
	if r != nil {
	}

	obj, _ := objectContainer.GetObject(context.Background(), filter.ByTypeOf[event.Multicaster]())

	multicaster := obj.(event.Multicaster)
	multicaster.MulticastEvent(nil, nil)

	event.ListenAsync(func(ctx context.Context, event event.Event) error {
		return nil
	})
}

type MyProcessor struct {
}

func newP(controller *UserController) *MyProcessor {
	return &MyProcessor{}
}

func (c MyProcessor) ProcessBeforeInit(ctx context.Context, object any) (any, error) {
	return object, nil
}

func (c MyProcessor) ProcessAfterInit(ctx context.Context, object any) (any, error) {
	return object, nil
}

type UserController struct {
	userService *UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: nil,
	}
}

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}
