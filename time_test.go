package iso8601_test

import (
	"encoding/json"
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/uudashr/iso8601"
)

type timePart struct {
	year    int
	month   time.Month
	day     int
	hour    int
	minute  int
	second  int
	nanosec int
	offset  int
}

func zoneOffset(t time.Time) int {
	_, offset := t.Zone()
	return offset
}

func nanoFrac(d int) int {
	l := len(strconv.Itoa(d)) - 1
	pow := math.Pow(10, float64(l))
	return d * 100000000 / int(pow)
}

func secOffset(hour, minute int) int {
	return (hour*60 + minute) * 60
}

func TestLayout(t *testing.T) {
	cases := map[string]struct {
		exp    string
		expect timePart
	}{
		"Use Z": {
			exp:    "2017-08-23T01:24:48.756Z",
			expect: timePart{2017, time.August, 23, 01, 24, 48, nanoFrac(756), 0},
		},

		"Use + sign": {
			exp:    "2017-08-23T01:24:48.756+07:00",
			expect: timePart{2017, time.August, 23, 01, 24, 48, nanoFrac(756), secOffset(7, 0)},
		},

		"Use - sign": {
			exp:    "2017-08-23T01:24:48.756-07:00",
			expect: timePart{2017, time.August, 23, 01, 24, 48, nanoFrac(756), secOffset(-7, 0)},
		},

		"No time fraction, with Z": {
			exp:    "2017-08-23T01:24:48Z",
			expect: timePart{2017, time.August, 23, 01, 24, 48, 0, 0},
		},

		"No time fraction, with + sign": {
			exp:    "2017-08-23T01:24:48+07:00",
			expect: timePart{2017, time.August, 23, 01, 24, 48, 0, secOffset(7, 0)},
		},

		"No time fraction, with - sign": {
			exp:    "2017-08-23T01:24:48-07:00",
			expect: timePart{2017, time.August, 23, 01, 24, 48, 0, secOffset(-7, 0)},
		},

		"1 digit time fraction": {
			exp:    "2017-08-23T01:24:48.1Z",
			expect: timePart{2017, time.August, 23, 01, 24, 48, nanoFrac(1), 0},
		},

		"2 digit time fraction": {
			exp:    "2017-08-23T01:24:48.27Z",
			expect: timePart{2017, time.August, 23, 01, 24, 48, nanoFrac(27), 0},
		},

		"5 digit time fraction": {
			exp:    "2017-08-23T01:24:48.87234Z",
			expect: timePart{2017, time.August, 23, 01, 24, 48, nanoFrac(87234), 0},
		},

		"9 digit time fraction": {
			exp:    "2017-08-23T01:24:48.987373613Z",
			expect: timePart{2017, time.August, 23, 01, 24, 48, nanoFrac(987373613), 0},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			parsed, err := time.Parse(iso8601.Layout, v.exp)
			if err != nil {
				t.Error("err:", err)
				return
			}

			if got, want := parsed.Year(), v.expect.year; got != want {
				t.Error("year got:", got, "want:", want)
			}

			if got, want := parsed.Month(), v.expect.month; got != want {
				t.Error("month got:", got, "want:", want)
			}

			if got, want := parsed.Day(), v.expect.day; got != want {
				t.Error("day got:", got, "want:", want)
			}

			if got, want := parsed.Hour(), v.expect.hour; got != want {
				t.Error("hour got:", got, "want:", want)
			}

			if got, want := parsed.Minute(), v.expect.minute; got != want {
				t.Error("minute got:", got, "want:", want)
			}

			if got, want := parsed.Second(), v.expect.second; got != want {
				t.Error("second got:", got, "want:", want)
			}

			if got, want := parsed.Nanosecond(), v.expect.nanosec; got != want {
				t.Error("nano got:", got, "want:", want)
			}

			if got, want := zoneOffset(parsed), v.expect.offset; got != want {
				t.Error("month got:", got, "want:", want)
			}
		})
	}
}

func TestLayout_Err(t *testing.T) {
	cases := map[string]string{
		"No time part": "2017-08-23",
		"No Z part":    "2017-08-23T01:24:48.756",
		"No date part": "01:24:48.756Z",
		"Empty":        "",
		"Space":        " ",
		"Double space": "  ",
		"Silly string": "silly string",
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			_, err := time.Parse(iso8601.Layout, v)
			if err == nil {
				t.Error("Expecting error for exp:", v)
			}
		})
	}
}

type body struct {
	OccuredOn iso8601.Time `json:"occuredOn"`
}

func TestUnmarshall(t *testing.T) {
	expect := timePart{2017, time.August, 23, 01, 24, 48, nanoFrac(756), secOffset(7, 0)}
	jsonBody := `{"occuredOn": "2017-08-23T01:24:48.756+07:00"}`
	var b body
	if err := json.Unmarshal([]byte(jsonBody), &b); err != nil {
		t.Error("err:", err)
	}

	parsed := time.Time(b.OccuredOn)
	if got, want := parsed.Year(), expect.year; got != want {
		t.Error("year got:", got, "want:", want)
	}

	if got, want := parsed.Month(), expect.month; got != want {
		t.Error("month got:", got, "want:", want)
	}

	if got, want := parsed.Day(), expect.day; got != want {
		t.Error("day got:", got, "want:", want)
	}

	if got, want := parsed.Hour(), expect.hour; got != want {
		t.Error("hour got:", got, "want:", want)
	}

	if got, want := parsed.Minute(), expect.minute; got != want {
		t.Error("minute got:", got, "want:", want)
	}

	if got, want := parsed.Second(), expect.second; got != want {
		t.Error("second got:", got, "want:", want)
	}

	if got, want := parsed.Nanosecond(), expect.nanosec; got != want {
		t.Error("nano got:", got, "want:", want)
	}

	if got, want := zoneOffset(parsed), expect.offset; got != want {
		t.Error("month got:", got, "want:", want)
	}
}

func TestUnmarshall_Error(t *testing.T) {
	cases := map[string]string{
		"Invalid date time format":      `{"occuredOn": "2017-08-23 01:24:48 +07:00"}`,
		"Empty date time":               `{"occuredOn": ""}`,
		"Empty with space on date time": `{"occuredOn": " "}`,
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			var b body
			err := json.Unmarshal([]byte(v), &b)
			if err == nil {
				t.Error("Expecting error for exp:", v)
			}
		})
	}
}

func TestMarshall(t *testing.T) {
	cases := map[string]struct {
		tick   time.Time
		expect string
	}{
		"UTC Time": {
			tick:   time.Date(2017, time.January, 15, 20, 11, 45, nanoFrac(245), time.UTC),
			expect: `{"occuredOn":"2017-01-15T20:11:45.245Z"}`,
		},

		"+00:00 offset": {
			tick:   time.Date(2017, time.January, 15, 20, 11, 45, nanoFrac(245), time.FixedZone("any", 0)),
			expect: `{"occuredOn":"2017-01-15T20:11:45.245Z"}`,
		},

		"+07:00 offset": {
			tick:   time.Date(2017, time.January, 15, 20, 11, 45, nanoFrac(245), time.FixedZone("any", secOffset(7, 0))),
			expect: `{"occuredOn":"2017-01-15T20:11:45.245+07:00"}`,
		},

		"-07:00 offset": {
			tick:   time.Date(2017, time.January, 15, 20, 11, 45, nanoFrac(245), time.FixedZone("any", secOffset(-7, 0))),
			expect: `{"occuredOn":"2017-01-15T20:11:45.245-07:00"}`,
		},
	}
	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			bd := body{OccuredOn: iso8601.Time(v.tick)}
			b, err := json.Marshal(bd)
			if err != nil {
				t.Error("err:", err)
				t.FailNow()
			}

			if got, want := string(b), v.expect; got != want {
				t.Error("got:", got, "want:", want)
			}
		})
	}
}
