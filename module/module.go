package module

import (
	"codnect.io/reflector"
	"fmt"
)

type Module interface {
	InitModule() error
}

func Use[M Module]() {
	typ := reflector.TypeOf[M]()
	if reflector.IsStruct(typ) {
		moduleStruct := reflector.ToStruct(typ)
		instance, err := moduleStruct.Instantiate()

		if err != nil {
			panic(err)
		}

		m := instance.Elem().(Module)
		err = m.InitModule()
		if err != nil {
			panic(fmt.Errorf("failed to initialize the module '%s': %e", moduleStruct.Name(), err))
		}
	}
}
