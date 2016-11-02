package main

import (
	"fmt"

	message "toby-cloud/go-client/message"
)

func main() {

	var tags = []string{"light"}

	tobymessage := message.NewMessage("hello", "TEXT", tags, "test")
	fmt.Println(tobymessage.String() + " yes!")
}
