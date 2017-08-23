package iso8601

import (
	"errors"
	"time"
)

// Layout of iso 8601 date time format.
const Layout = "2006-01-02T15:04:05.999999999Z07:00"

// Time with json marshal and unmarshal capability of iso 8601 format.
type Time time.Time

// UnmarshalJSON implements the json.Unmarshaller interface.
func (jt *Time) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}

	parsed, err := time.Parse(`"`+Layout+`"`, string(b))
	if err != nil {
		return err
	}

	*jt = Time(parsed)
	return nil
}

// MarshalJSON implements the json.Marshaller interface.
func (jt Time) MarshalJSON() ([]byte, error) {
	t := time.Time(jt)
	if y := t.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("iso8601: year outside of range [0,9999]")
	}
	b := make([]byte, 0, len(Layout)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, Layout)
	b = append(b, '"')
	return b, nil
}
