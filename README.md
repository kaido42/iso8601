# ISO 8601

JSON time serialization support [ISO 8601](https://xml2rfc.tools.ietf.org/public/rfc/html/rfc3339.html#anchor14) specification.

## Installation
`go get github.com/uudashr/iso8601`

## Usage
Use the time Layout (same with time.RFC3339Nano)

```golang
func parseTime(s string) (t time.Time, err error) {
    t, err = time.Parse(iso8601.Layout, s)
    return
}
```
