package main

import "fmt"


func main() {
	depots, customers, err := LoadProblem("problems/p01")
	if err != nil {
		panic(err)
	}
	
	fmt.Println(customers)
	fmt.Println(depots)
}

