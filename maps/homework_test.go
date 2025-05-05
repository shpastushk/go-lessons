package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type node struct {
	key, value  int
	left, right *node
}

func (n *node) processNode(action func(int, int)) {
	if n == nil {
		return
	}

	n.left.processNode(action)
	action(n.key, n.value)
	n.right.processNode(action)
}

type OrderedMap struct {
	root *node
	size int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{}
}

func (m *OrderedMap) Insert(key, value int) {
	x := m.root
	var y *node
	for x != nil {
		if x.key == key {
			x.value = value
			return
		}

		y = x
		if x.key > key {
			x = x.left
		} else {
			x = x.right
		}
	}

	newNode := &node{
		key:   key,
		value: value,
	}

	if y == nil {
		m.root = newNode
		m.size++
		return
	}

	if y.key > key {
		y.left = newNode
	} else {
		y.right = newNode
	}
	m.size++

}

func (m *OrderedMap) Erase(key int) {
	x := m.root
	var y *node
	for x != nil {
		if x.key == key {
			break
		} else {
			y = x
			if x.key > key {
				x = x.left
			} else {
				x = x.right
			}
		}
	}

	if x == nil {
		return
	}

	if x.right == nil {
		if y == nil {
			m.root = x.left
		} else {
			if x == y.left {
				y.left = x.left
			} else {
				y.right = x.left
			}
		}
		m.size--
		return
	}

	leftMost := x.right
	y = nil
	for leftMost.left != nil {
		y = leftMost
		leftMost = leftMost.left
	}
	if y != nil {
		y.left = leftMost.right
	} else {
		x.right = leftMost.right
	}
	x.key = leftMost.key
	x.value = leftMost.value
	m.size--
}

func (m *OrderedMap) Contains(key int) bool {
	x := m.root
	for x != nil {
		if x.key == key {
			return true
		}
		if x.key > key {
			x = x.left
		} else {
			x = x.right
		}
	}
	return false
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	m.root.processNode(action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
