[![Build Status](https://travis-ci.org/uudashr/iso8601.svg?branch=master)](https://travis-ci.org/uudashr/iso8601)
[![GoDoc](https://godoc.org/github.com/uudashr/iso8601?status.svg)](https://godoc.org/github.com/uudashr/iso8601)
# ISO 8601

JSON time serialization support [ISO 8601](https://xml2rfc.tools.ietf.org/public/rfc/html/rfc3339.html#anchor14) specification.

## Installation
`go get github.com/uudashr/iso8601`

## Usage
Use the time Layout (same with time.RFC3339Nano)

```golang
import "time"

func parseTime(s string) (t time.Time, err error) {
    t, err = time.Parse(iso8601.Layout, s)
    return
}
```

Use it on struct
```golang
import (
    "fmt"
    "json"
    "time"
)

type Event struct {
    Name string `json:"name"`
    OccuredOn iso8601.Time `json:"occuredOn"`
}

now := time.Now()
event := Event {Name: "Sign In", iso8601.Time(Time:now)}
b, _ := json.Marshal(event)    
fmt.Println(string(b)) // show the marshalled struct
```

Unmarshal into struct
```golang
import (
)

type Event struct {
    Name string `json:"name"`
    OccuredOn iso8601.Time `json:"occuredOn"`
}

source := "{\"name\":\"test\",\"occuredOn\":\"2002-10-02T10:00:00-05:00\"}"
var event Event
unMarshal(source, &event)

fmt.Prinln("occured at unixtime: %d", event.OccuredOn.Unix())

```
