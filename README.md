# toby-go

> Use Toby with any Go application.

## Installation

```golang
import (
	tobyBot "github.com/toby-cloud/toby-go/bot"
	tobyMessage "github.com/toby-cloud/toby-go/message"
)
```

Then do `go get` to install the Toby package.


## Usage

See the `examples` folder.


## Testing

Use `go test` to run unit tests. The MQTT connection is mocked.


### TODO:

* Add unit testing to `message.go`
* Add unit testing to `hashtag.go`
* `utils/hashtag.go` is unused
* Update examples
