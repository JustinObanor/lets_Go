package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	s := []string{"this is my is is my wow so is is my my so hey no yes"}
	sc := rangeMap(s, 2)
	for k, v := range sc {
		fmt.Printf("%s\t appeared : %d\n", k, v)
	}
	//is : 5
	//my : 4
	//so : 2

}

func rangeMap(s []string, count int) map[string]int {

	words := strings.Split(strings.Join(s, ""), " ")
	m := make(map[string]int)

	for _, word := range words {
		_, ok := m[word]
		if ok {
			m[word]++
		} else {
			m[word] = 1
		}
	}

	counts := make(map[string]int)
	for key, value := range m {
		if value > 1 {
			//copying the value from the main map to the new one
			counts[key] = value
		}
	}

	keys := make([]string, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i int, j int) bool {
		return counts[keys[i]] > counts[keys[j]]
	})

	// Builds result map
	result := make(map[string]int)
	for _, key := range keys {
		result[key] = counts[key]
		count--
		if count == 0 {
			break
		}

	}
	return result
}
