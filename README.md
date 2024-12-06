## Install

```bash
go get github.com/ihatiko/environ@latest
```

### Description

todo

### Api


```go
package main

import (
	"os"

	"github.com/ihatiko/environ"
)

type A struct {
	Field1 string
}

func main() {
	os.Setenv("Field1", "customvalue")
	result := new(A)
	environ.Parse(result) // after parse result.Field1 equals customvalue
}
```


## Supported

```go
url.Url

time.Duration

time.Time (format "2006-01-02 15:04:05")
```