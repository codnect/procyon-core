package main

import (
	core "codnect.io/procyon-core"
	"codnect.io/procyon-core/component"
	"codnect.io/procyon-core/module"
	"codnect.io/procyon-core/procyon"
	"codnect.io/procyon-core/runtime"
	"context"
)

func init() {
	component.Register(newMyCommandLineApp)
	module.Use[core.Module]()
}

func main() {
	err := procyon.New().Run()
	if err != nil {
		panic(err)
	}
}

type MyCommandLineApp struct {
}

func newMyCommandLineApp() MyCommandLineApp {
	return MyCommandLineApp{}
}

func (a MyCommandLineApp) Run(ctx context.Context, args *runtime.Arguments) error {

	return nil
}

func canConvert(obj any) {

}
