package main

import (
	"fmt"

	message "github.com/toby-cloud/toby-go/message"
)

func main() {

	var tags = []string{"light"}

	tobymessage := message.NewMessage("hello", "TEXT", tags, "test")
	fmt.Println(tobymessage.String() + " yes!")
}
