package main

import (
	"fmt"
)

//interface{}
func mapKeyReader(stringMap map[string]string) {

	fmt.Println("Map values: ____________________________________________________")
	for k := range stringMap {
		fmt.Println(k)
	}
	fmt.Println("Map end: _______________________________________________________")
}
