//go:build !cgo || (!linux && !windows)

package localeinfo

import (
	"time"
)

type stubLocale struct {
}

func NewLocale(name string) (*stubLocale, error) {
	return &stubLocale{}, nil
}

func (l *stubLocale) Encoding() string {
	return ""
}

func (l *stubLocale) DateTimeFormat() string {
	return ""
}

func (l *stubLocale) DateFormat() string {
	return ""
}

func (l *stubLocale) TimeFormat() string {
	return ""
}

func (l *stubLocale) AM() string {
	return ""
}

func (l *stubLocale) PM() string {
	return ""
}

func (l *stubLocale) TimeAMPMFormat() string {
	return ""
}

func (l *stubLocale) Day(day time.Weekday) string {
	return ""
}

func (l *stubLocale) ShortDay(day time.Weekday) string {
	return ""
}

func (l *stubLocale) Month(month time.Month) string {
	return ""
}

func (l *stubLocale) ShortMonth(month time.Month) string {
	return ""
}

func (l *stubLocale) Radix() string {
	return ""
}

func (l *stubLocale) ThousandSeparator() string {
	return ""
}

func (l *stubLocale) Currency() (symbol string, position CurrencyFormat) {
	return
}
