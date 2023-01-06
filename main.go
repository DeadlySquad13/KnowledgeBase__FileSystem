package main

import (
	"fmt"

	gremlingo "github.com/apache/tinkerpop/gremlin-go/v3/driver"
)

// syntactic sugar
var __ = gremlingo.T__
var gt = gremlingo.P.Gt
var order = gremlingo.Order


type FileSystem struct {
    driverRemoteConnection gremlingo.DriverRemoteConnection
    g *gremlingo.GraphTraversalSource
}
type Bookmark struct {
    path, name string
}

func NewFileSystem(connectionString string) (*FileSystem, error) {
    // Creating the connection to the server.
	driverRemoteConnection, err := gremlingo.NewDriverRemoteConnection(connectionString)
	if err != nil {
		return nil, err
	}
	// Creating graph traversal
	g := gremlingo.Traversal_().WithRemote(driverRemoteConnection)

    fs := &FileSystem{driverRemoteConnection: *driverRemoteConnection, g: g}

    return fs, nil
}


func (fs *FileSystem) addBookmark(b Bookmark) <-chan error {
    promise := fs.g.AddV("bookmark").Property("path", b.path).Property("name", b.name).Iterate()

    return promise
}

// FIX: Rise errors to stop upper level executions if something wasn't added.
func (fs *FileSystem) addBookmarks(bookmarks []Bookmark) {
    for _, bookmark := range bookmarks {
        // Wait for all steps to finish execution and check for error.
        promiseErr := <-fs.addBookmark(bookmark)
        if promiseErr != nil {
            fmt.Println(promiseErr)
            return 
        }
    }
}

// Perform traversal
func (fs *FileSystem) getBookmarks() []*gremlingo.Result {
    result, err := fs.g.V().HasLabel("bookmark").Values("path").ToList()
	if result == nil {
		return nil
	}
	if err != nil {
        fmt.Println(err)
		return nil
	}

	return result
}

func (fs *FileSystem) getBookmarksByName(name string) []*gremlingo.Result {
    result, err := fs.g.V().HasLabel("bookmark").Has("name", name).Values("path").ToList()
	if result == nil {
		return nil
	}
	if err != nil {
        fmt.Println(err)
		return nil
	}

	return result
}

func main() {
    fs, err := NewFileSystem("ws://localhost:8182/gremlin")
    if err != nil {
        fmt.Println(err)
        return
    }

    // fs.addBookmarks([]Bookmark{
    //     { name: "FileSystem", path: "S:\\ds13\\Projects\\--personal\\KnowledgeBase\\FileSystem\\" },
    //     { name: "Bomonka", path: "E:\\Projects\\--educational\\Bomonka\\" },
    //     { name: "Configs", path: "E:\\Scripts\\" },
    // })

    vertices := fs.getBookmarksByName("Bomonka")
    for _, v := range vertices {
		fmt.Println(v.GetString())
	}

    // Cleanup
    defer fs.driverRemoteConnection.Close()
}
