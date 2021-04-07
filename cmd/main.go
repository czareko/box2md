package main

import (
	"devopsbox.io/box2md/pkg/box2md"
	"fmt"
	"os"
)

func main() {
	schemaJson := os.Args[1]

	schema := box2md.Read(schemaJson)

	fmt.Println(schema.ToMDString(0))
}
