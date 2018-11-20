package episode

import (
	"fmt"
	"time"
)

type Episode struct {
	Id      int
	Name    string
	Number  int
	Season  int
	AirDate string
	Url     string
}

func (e Episode) diff() (year, month, day int) {
	layout := "2006-01-02"
	date, err := time.Parse(layout, e.AirDate)

	if err != nil {
		fmt.Println(err)
	}

	return diff(date, time.Now())
}

func (e Episode) DaysLeft() int {
	year, month, day := e.diff()

	if month > 0 {
		day += month / 30
	}

	if year > 0 {
		day += year / 360
	}

	return day
}

func (e Episode) TimeLeft() string {
	// _, _, day, _, _, _ := e.diff()
	year, month, day := e.diff()

	if year > 0 {
		return fmt.Sprintf("%d Year/s, %d Month/s and %d Day/s", year, month, day)
	} else if month > 0 {
		return fmt.Sprintf("%d Month/s and %d Day/s", month, day)
	} else {
		return fmt.Sprintf("%d Day/s", day)
	}
}

func diff(a, b time.Time) (year, month, day int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)

	// Normalize negative values
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}
