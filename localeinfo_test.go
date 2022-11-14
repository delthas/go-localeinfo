package localeinfo

import (
	"testing"
	"time"
)

func assert(t *testing.T, got string, want string, fmt string) {
	if want == got {
		return
	}
	t.Errorf(fmt+": want: %v: got: %v", want, got)
}

func TestDateFmt(t *testing.T) {
	l, err := NewLocale("")
	if err != nil {
		t.Errorf("NewLocale: %v", err)
	}
	assert(t, l.DateFormat(), "%m/%d/%y", "DateFormat")
	assert(t, l.DateTimeFormat(), "%a %b %e %H:%M:%S %Y", "DateTimeFormat")
	assert(t, l.Month(time.February), "February", "Month")
	assert(t, l.Day(time.Friday), "Friday", "Day")
	assert(t, l.AM(), "AM", "AM")
	assert(t, l.PM(), "PM", "PM")
	s, _ := l.Currency()
	assert(t, s, "", "Currency")
	assert(t, l.Encoding(), "UTF-8", "Encoding")
	assert(t, l.Radix(), ".", "Radix")
	assert(t, l.ThousandSeparator(), "", "ThousandSeparator")
	assert(t, l.TimeFormat(), "%H:%M:%S", "TimeFormat")
	assert(t, l.TimeAMPMFormat(), "%I:%M:%S %p", "TimeAMPMFormat")
	assert(t, l.ShortDay(time.Wednesday), "Wed", "ShortDay")
	assert(t, l.ShortMonth(time.July), "Jul", "ShortMonth")
}
