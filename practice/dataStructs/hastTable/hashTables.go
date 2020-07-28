package main

type hashTable struct {
	data map[(interface{})][]interface{}
}

func (h *hashTable) get(key interface{}) interface{} {
	return h.data[key]
}

func (h *hashTable) set(key, value interface{}) {
	v, ok := h.data[key]

	if !ok {
		v = make([]interface{}, 0)
		v = append(v, value)
	}

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
