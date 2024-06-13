package main

import (
	"JamalMartin/golang-dev-logic-challenge-MartinJamal/routes"
	"fmt"
)

func main() {
	r := routes.SetupRouter()
	if err := r.Run(":8080"); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}

}
