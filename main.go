package main

import (
	"fmt"

	gremlingo "github.com/apache/tinkerpop/gremlin-go/v3/driver"
)

// syntactic sugar
var __ = gremlingo.T__
var gt = gremlingo.P.Gt
var order = gremlingo.Order


func addBookmark(g *gremlingo.GraphTraversalSource, path string) <-chan error {
    promise := g.AddV("bookmark").Property("path", path).Iterate()

    return promise
}

func main() {
	// Creating the connection to the server.
	driverRemoteConnection, err := gremlingo.NewDriverRemoteConnection("ws://localhost:8182/gremlin")
	if err != nil {
		fmt.Println(err)
		return
	}
	// Creating graph traversal
	g := gremlingo.Traversal_().WithRemote(driverRemoteConnection)

    
    // Wait for all steps to finish execution and check for error.
    promiseErr := <-addBookmark(g, "S:\\ds13\\Projects\\--personal\\KnowledgeBase\\FileSystem")
    if promiseErr != nil {
        fmt.Println(promiseErr)
        return
    }

	// Perform traversal
    result, err := g.V().Values().ToList()
	// result, err := g.V().HasLabel("person").Has("age", __.Is(gt(28))).Order().By("age", order.Desc).Values("name").ToList()
	if err != nil {
		fmt.Println(err)
		return
	}
    fmt.Println("test")
	for _, r := range result {
		fmt.Println(r.GetString())
	}

    // Cleanup
    defer driverRemoteConnection.Close()
}
