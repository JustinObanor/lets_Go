package main

import "fmt"

type hashTable struct {
	data map[(interface{})][]interface{}
}

func New() hashTable {
	return hashTable{
		data: make(map[interface{}][]interface{}),
	}
}

func (h *hashTable) get(key interface{}) interface{} {
	return h.data[key]
}

func (h *hashTable) set(key, value interface{}) {
	v := h.data[key]

	v = append(v, value)

	h.data[key] = v
}

func (h *hashTable) keys() []interface{} {
	res := make([]interface{}, 0)
	for key := range h.data {
		res = append(res, key)
	}
	return res
}

func (h *hashTable) values() []interface{} {
	res := make([]interface{}, 0)
	for _, value := range h.data {
		res = append(res, value)
	}
	return res
}

func main() {
	h := New()
	h.set("a", 1)
	h.set("a", 2)
	fmt.Println(h.get("a"))
}
