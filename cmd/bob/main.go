package main

import (
	"github.com/gen2brain/beeep"
)

func main() {
	err := beeep.Notify("Hi", "There", "")
	if err != nil {
		panic(err)
	}
}
