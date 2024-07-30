package component

import (
	"codnect.io/reflector"
	"context"
	"fmt"
	"github.com/codnect/procyoncore/component/filter"
	"sync"
)

type SingletonRegistry interface {
	Register(name string, object any) error
	Remove(name string) error
	Find(filters ...filter.Filter) (any, error)
	FindFirst(filters ...filter.Filter) (any, bool)
	List(filters ...filter.Filter) []any
	OrElseCreate(name string, provider ObjectProvider) (any, error)
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
	return &SingletonObjectRegistry{
		singletonObjects:              make(map[string]any),
		singletonObjectsInPreparation: make(map[string]struct{}),
		typesOfSingletonObjects:       make(map[string]reflector.Type),
	}
}

func (r *SingletonObjectRegistry) Register(name string, object any) error {
	defer r.muSingletonObjects.Unlock()
	r.muSingletonObjects.Lock()

	if _, exists := r.singletonObjects[name]; exists {
		return fmt.Errorf("object with name %s already exists", name)
	}

	r.singletonObjects[name] = object
	r.typesOfSingletonObjects[name] = reflector.TypeOfAny(object)
	return nil
}

func (r *SingletonObjectRegistry) Remove(name string) error {
	defer r.muSingletonObjects.Unlock()
	r.muSingletonObjects.Lock()

	if _, exists := r.singletonObjects[name]; !exists {
		return fmt.Errorf("no found object with name %s", name)
	}

	delete(r.singletonObjects, name)
	delete(r.typesOfSingletonObjects, name)
	return nil
}

func (r *SingletonObjectRegistry) Find(filters ...filter.Filter) (any, error) {
	objectList := r.List(filters...)

	if len(objectList) > 1 {
		return nil, fmt.Errorf("objects cannot be distinguished because too many matching found")
	}

	if len(objectList) == 0 {
		filterOpts := filter.Of(filters...)

		return nil, ObjectNotFoundError{
			name: filterOpts.Name,
			typ:  filterOpts.Type,
		}
	}

	return objectList[0], nil
}

func (r *SingletonObjectRegistry) FindFirst(filters ...filter.Filter) (any, bool) {
	objectList := r.List(filters...)

	if len(objectList) == 0 {
		return nil, false
	}

	return objectList[0], true
}

func (r *SingletonObjectRegistry) List(filters ...filter.Filter) []any {
	defer r.muSingletonObjects.Unlock()
	r.muSingletonObjects.Lock()

	filterOpts := filter.Of(filters...)
	objectList := make([]any, 0)

	for objectName, objectType := range r.typesOfSingletonObjects {

		if filterOpts.Name != "" && filterOpts.Name != objectName {
			continue
		}

		if filterOpts.Type == nil {
			objectList = append(objectList, r.singletonObjects[objectName])
			continue
		}

		if objectType.CanConvert(filterOpts.Type) {
			objectList = append(objectList, r.singletonObjects[objectName])
		} else if reflector.IsPointer(objectType) && !reflector.IsPointer(filterOpts.Type) && !reflector.IsInterface(filterOpts.Type) {
			ptrType := reflector.ToPointer(objectType)

			if ptrType.Elem().CanConvert(filterOpts.Type) {
				val, err := ptrType.Elem().Value()

				if err == nil {
					objectList = append(objectList, val)
				}
			}
		}

	}

	return objectList
}

func (r *SingletonObjectRegistry) OrElseCreate(name string, provider ObjectProvider) (any, error) {
	object, err := r.Find(filter.ByName(name))

	if err == nil {
		return object, nil
	}

	err = r.putObjectToPreparation(name)

	if err != nil {
		return nil, err
	}

	defer r.removeObjectFromPreparation(name)

	object, err = provider(context.Background())

	if err != nil {
		return nil, err
	}

	r.singletonObjects[name] = object
	r.typesOfSingletonObjects[name] = reflector.TypeOfAny(object)

	return object, nil
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
		return fmt.Errorf("object with name %s is currently in preparation, maybe it has got circular dependency cycle", name)
	}

	r.singletonObjectsInPreparation[name] = struct{}{}
	return nil
}

func (r *SingletonObjectRegistry) removeObjectFromPreparation(name string) {
	defer r.muSingletonObjects.Unlock()
	r.muSingletonObjects.Lock()
	delete(r.singletonObjectsInPreparation, name)
}
