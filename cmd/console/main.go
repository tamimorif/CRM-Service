package main

import (
	"fmt"
	"github.com/google/uuid"
	"log"
)

func main() {
	id := uuid.New()
	fmt.Println(id)
	id, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
	fmt.Printf("%T %v", id, id)
}
