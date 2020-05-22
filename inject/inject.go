package inject

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrInvalidValue           = errors.New("invalid object")
	ErrCallbackMustBeFunction = errors.New("callback must be a function")
	ErrValueNotExists         = errors.New("value not exists")
	defaultContainer          = New()
)

type Container struct {
	m         map[string]interface{}
	callbacks []reflect.Value
}

func New() *Container {
	return &Container{
		m:         make(map[string]interface{}),
		callbacks: make([]reflect.Value, 0),
	}
}

// Provide 注入
func (c *Container) Provide(object interface{}) {
	identifier := getFullIdentifier(reflect.TypeOf(object))
	c.m[identifier] = object
}

// AddCallback 添加回调函数，调用Populate后会自动执行
func (c *Container) AddCallback(callback interface{}) error {
	value := reflect.ValueOf(callback)
	if value.Kind() != reflect.Func || value.IsZero() {
		return ErrCallbackMustBeFunction
	}
	c.callbacks = append(c.callbacks, value)
	return nil
}

// Populate 注入
func (c *Container) Populate() error {
	// 注入
	for _, object := range c.m {
		// 只往结构体的字段注入
		objectValue := reflect.ValueOf(object)
		for !objectValue.IsZero() && objectValue.Kind() == reflect.Ptr {
			objectValue = objectValue.Elem()
		}
		if objectValue.Kind() != reflect.Struct {
			continue
		}

		// 往结构体的字段注入值
		numField := objectValue.NumField()
		for i := 0; i < numField; i++ {
			field := objectValue.Field(i)
			identifier := getFullIdentifier(field.Type())

			// 是否存在(跟结构体字段)相同类型的依赖
			object, ok := c.m[identifier]
			if !ok {
				continue
			}

			// 是否能赋值
			fieldValue := reflect.ValueOf(object)
			if !field.Type().AssignableTo(fieldValue.Type()) || !field.CanSet() {
				continue
			}
			field.Set(fieldValue)
		}
	}

	// 执行回调
	for _, callback := range c.callbacks {
		callbackType := callback.Type()

		// 构造函数调用参数
		var inParams []reflect.Value
		numIn := callbackType.NumIn()
		for i := 0; i < numIn; i++ {
			paramType := callbackType.In(i)
			identifier := getFullIdentifier(paramType)

			// 是否存在(跟函数参数)相同类型的依赖
			object, ok := c.m[identifier]
			if !ok {
				return errors.New(fmt.Sprintf("value of %s is required for callbacks", identifier))
			}
			inParams = append(inParams, reflect.ValueOf(object))
		}
		callback.Call(inParams)
	}
	return nil
}

// Get 获取指定类型值，应该在Populate后调用
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

// Get 获取指定类型值，应该在Populate后调用
func Get(pValues ...interface{}) error {
	return defaultContainer.Get(pValues...)
}

// AddCallback 添加回调函数，调用Populate后会自动执行
func AddCallback(callback interface{}) error {
	return defaultContainer.AddCallback(callback)
}
