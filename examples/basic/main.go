package main

import (
	"fmt"

	hapesay "github.com/hapefarts/hapefarts/v2"
)

func main() {
	if false {
		simple()
	} else {
		complex()
	}
}

func simple() {
	say, err := hapesay.Say(
		"Hello",
		hapesay.Type("mobile"),
		hapesay.BallonWidth(40),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(say)
}

func complex() {
	hape, err := hapesay.New(
		hapesay.BallonWidth(40),
		//hapesay.Thinking(),
		hapesay.Random(),
	)
	if err != nil {
		panic(err)
	}
	say, err := hape.Say("Hello")
	if err != nil {
		panic(err)
	}
	fmt.Println(say)
}
