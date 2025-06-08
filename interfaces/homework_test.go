package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}
type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}

type Container struct {
	// need to implement
	deps map[string]interface{}
}

func NewContainer() *Container {
	// need to implement
	return &Container{deps: make(map[string]interface{})}
}

func (c *Container) RegisterType(name string, constructor interface{}) {
	// need to implement
	c.deps[name] = constructor
}

func (c *Container) Resolve(name string) (interface{}, error) {
	// need to implement
	f := c.deps[name]
	fn := reflect.ValueOf(f)
	if reflect.ValueOf(f).Kind() != reflect.Func {
		return nil, fmt.Errorf("not a function")
	}

	results := fn.Call([]reflect.Value{})
	value := results[0]

	switch name {
	case "UserService":
		return value.Interface().(*UserService), nil
	case "MessageService":
		return value.Interface().(*MessageService), nil
	default:
		return nil, fmt.Errorf("unknown type: %s", name)
	}
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)
}
