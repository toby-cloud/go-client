package main

import (
	"./toby/message"
	"fmt"
)


func main() {

	var tags = []string{"light"}

	tobymessage := message.NewMessage("hello", "TEXT", tags, "test")
	fmt.Println(tobymessage.String() + " yes!")
}
