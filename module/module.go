package module

import "codnect.io/reflector"

type Module interface {
	InitModule()
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
		m.InitModule()
	}
}
