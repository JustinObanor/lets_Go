package main

import "fmt"

type Graph struct {
	numOfNodes   int
	adjacentList map[int][]int
}

func newGraph() Graph {
	return Graph{
		numOfNodes:   0,
		adjacentList: make(map[int][]int),
	}
}

func (g *Graph) addVertex(node int) {
	g.adjacentList[node] = make([]int, 0, 3)
	g.numOfNodes++
}

func (g *Graph) addEdge(node1, node2 int) {
	g.adjacentList[node1] = append(g.adjacentList[node1], node2)
	g.adjacentList[node2] = append(g.adjacentList[node2], node1)

}

func (g *Graph) showConnections() {

}

func main() {
	g := newGraph()
	for i := 0; i <= 6; i++ {
		g.addVertex(i)
	}

	g.addEdge(1, 0)
	g.addEdge(1, 2)
	g.addEdge(3, 1)
	g.addEdge(3, 4)
	g.addEdge(4, 2)
	g.addEdge(4, 5)
	g.addEdge(0, 2)
	g.addEdge(6, 5)

	keys := make([]int, 0, g.numOfNodes)
	for k := range g.adjacentList {
		keys = append(keys, k)
	}

	for k := range keys {
		fmt.Printf("%d --> %v\n", k, g.adjacentList[k])
	}
}
