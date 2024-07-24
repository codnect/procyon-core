package component

import (
	"codnect.io/reflector"
	"context"
	"errors"
	"fmt"
	"sync"
)

type SingletonRegistry interface {
	Register(name string, instance any) error
	Find(name string) (any, bool)
	FindByType(requiredType reflector.Type) (any, error)
	ListByType(requiredType reflector.Type) []any
	OrElseGet(name string, provider ObjectProvider) (any, error)
	Contains(name string) bool
	Names() []string
	Count() int
}

type SingletonObjectRegistry struct {
	singletonObjects              map[string]any
	singletonObjectsInPreparation map[string]struct{}
	typesOfSingletonObjects       map[string]reflector.Type
	muSingletonObjects            sync.RWMutex
}

func NewSingletonObjectRegistry() *SingletonObjectRegistry {
	return &SingletonObjectRegistry{}
}

func (r *SingletonObjectRegistry) Register(name string, instance any) error {
	defer r.muSingletonObjects.Unlock()
	r.muSingletonObjects.Lock()

	if _, exists := r.singletonObjects[name]; exists {
		return fmt.Errorf("instance with name %s already exists", name)
	}

	r.singletonObjects[name] = instance
	r.typesOfSingletonObjects[name] = reflector.TypeOfAny(instance)
	return nil
}

func (r *SingletonObjectRegistry) Find(name string) (any, bool) {
	defer r.muSingletonObjects.Unlock()
	r.muSingletonObjects.Lock()

	if instance, exists := r.singletonObjects[name]; exists {
		return instance, true
	}

	return nil, false
}

func (r *SingletonObjectRegistry) FindByType(requiredType reflector.Type) (any, error) {
	if requiredType == nil {
		return nil, errors.New("container: requiredType cannot be nil")
	}

	instances := r.ListByType(requiredType)
	if len(instances) > 1 {
		return nil, fmt.Errorf("container: instances cannot be distinguished for required type %s", requiredType.Name())
	}

	if len(instances) == 0 {
		return nil, &NotFoundError{
			//ErrorString: fmt.Sprintf("container: not found any instance of type %s", requiredType.Name()),
		}
	}

	return instances[0], nil
}

func (r *SingletonObjectRegistry) ListByType(requiredType reflector.Type) []any {
	defer r.muSingletonObjects.Unlock()
	r.muSingletonObjects.Lock()

	instances := make([]any, 0)

	for name, typ := range r.typesOfSingletonObjects {

		if typ.CanConvert(requiredType) {
			instances = append(instances, r.singletonObjects[name])
		} else if reflector.IsPointer(typ) && !reflector.IsPointer(requiredType) && !reflector.IsInterface(requiredType) {
			ptrType := reflector.ToPointer(typ)

			if ptrType.Elem().CanConvert(requiredType) {
				val, err := ptrType.Elem().Value()

				if err == nil {
					instances = append(instances, val)
				}
			}
		}

	}

	return instances
}

func (r *SingletonObjectRegistry) OrElseGet(name string, provider ObjectProvider) (any, error) {
	instance, ok := r.Find(name)

	if ok {
		return instance, nil
	}

	err := r.putObjectToPreparation(name)

	if err != nil {
		return nil, err
	}

	defer func() {
		r.removeObjectFromPreparation(name)
	}()

	instance, err = provider(context.Background())

	if err != nil {
		return nil, err
	}

	r.singletonObjects[name] = instance
	r.typesOfSingletonObjects[name] = reflector.TypeOfAny(instance)

	return instance, nil
}

func (r *SingletonObjectRegistry) Contains(name string) bool {
	defer r.muSingletonObjects.Unlock()
	r.muSingletonObjects.Lock()

	_, exists := r.singletonObjects[name]
	return exists
}

func (r *SingletonObjectRegistry) Names() []string {
	defer r.muSingletonObjects.Unlock()
	r.muSingletonObjects.Lock()

	names := make([]string, 0)
	for name := range r.singletonObjects {
		names = append(names, name)
	}

	return names
}

func (r *SingletonObjectRegistry) Count() int {
	defer r.muSingletonObjects.Unlock()
	r.muSingletonObjects.Lock()

	return len(r.singletonObjects)
}

func (r *SingletonObjectRegistry) putObjectToPreparation(name string) error {
	defer r.muSingletonObjects.Unlock()
	r.muSingletonObjects.Lock()

	if _, ok := r.singletonObjectsInPreparation[name]; ok {
		return fmt.Errorf("instance with name %s is currently in preparation, maybe it has got circular dependency cycle", name)
	}

	r.singletonObjectsInPreparation[name] = struct{}{}
	return nil
}

func (r *SingletonObjectRegistry) removeObjectFromPreparation(name string) {
	defer r.muSingletonObjects.Unlock()
	r.muSingletonObjects.Lock()
	delete(r.singletonObjectsInPreparation, name)
}
