package inject

import (
	"errors"
	"reflect"
)

var (
	ErrInvalidValue   = errors.New("invalid object")
	ErrValueNotExists = errors.New("value not exists")
	defaultContainer  = New()
)

type Container struct {
	m map[string]interface{}
}

func New() *Container {
	return &Container{m: make(map[string]interface{})}
}

// Provide 注入
func (c *Container) Provide(object interface{}) {
	identifier := getFullIdentifier(reflect.TypeOf(object))
	c.m[identifier] = object
}

// Populate 注入
func (c *Container) Populate() error {
	for _, object := range c.m {
		objectValue := reflect.ValueOf(object)
		for !objectValue.IsZero() && objectValue.Kind() == reflect.Ptr {
			objectValue = objectValue.Elem()
		}
		if objectValue.Kind() != reflect.Struct {
			continue
		}

		numField := objectValue.NumField()
		for i := 0; i < numField; i++ {
			field := objectValue.Field(i)
			identifier := getFullIdentifier(field.Type())

			object, ok := c.m[identifier]
			if !ok {
				continue
			}
			fieldValue := reflect.ValueOf(object)

			if !field.Type().AssignableTo(fieldValue.Type()) || !field.CanSet() {
				continue
			}
			field.Set(fieldValue)
		}
	}
	return nil
}

// Get 获取指定类型值
func (c *Container) Get(pValues ...interface{}) error {
	for _, pValue := range pValues {
		value := reflect.ValueOf(pValue)
		if value.IsZero() || value.Kind() != reflect.Ptr {
			return ErrInvalidValue
		}

		elemValue := value.Elem()
		identifier := getFullIdentifier(elemValue.Type())
		object, ok := c.m[identifier]
		if !ok {
			return ErrValueNotExists
		}

		objectValue := reflect.ValueOf(object)
		if !elemValue.Type().AssignableTo(objectValue.Type()) {
			return ErrValueNotExists
		}
		elemValue.Set(objectValue)
	}
	return nil
}

// Provide 注入
func Provide(object interface{}) {
	defaultContainer.Provide(object)
}

// Populate 注入
func Populate() error {
	return defaultContainer.Populate()
}

// Get 获取指定类型值
func Get(pValues ...interface{}) error {
	return defaultContainer.Get(pValues...)
}
