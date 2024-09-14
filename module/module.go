package module

import (
	"fmt"
	"reflect"
)

type Module interface {
	InitModule() error
}

func Use[M Module]() {
	moduleType := reflect.TypeFor[M]()
	if moduleType.Kind() == reflect.Struct {
		moduleValue := reflect.New(moduleType)

		m := moduleValue.Interface().(Module)
		err := m.InitModule()
		if err != nil {
			panic(fmt.Errorf("failed to initialize the module '%s': %e", moduleType.Name(), err))
		}
	}
}
