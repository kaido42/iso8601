package iso8601_test

import (
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
