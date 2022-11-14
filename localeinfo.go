// Package localeinfo offers platform-agnostic information about the current locale.
// The main entrypoints are Default and NewLocale, which return a Locale.
// Most calls of interest are the methods of Locale.
package localeinfo

import (
	"time"
)

// Locale represents a particular locale, instantiated from Default or NewLocale.
// Locale is an interface over platform-specific locale implementations.
// Most calls are best-effort, as not all platforms support the same locale information.
type Locale interface {
	// Encoding returns the name of the character encoding used in the selected locale, such as "UTF-8", "ISO-8859-1", or "ANSI_X3.4-1968" (better known as US-ASCII).
	Encoding() string

	// DateTimeFormat returns a format string to represent time and date in a locale-specific way (corresponds to strftime's %c).
	DateTimeFormat() string

	// DateFormat returns a format string to represent a date in a locale-specific way (corresponds to strftime's %x).
	DateFormat() string

	// TimeFormat returns a format string to represent a time in a locale-specific way (corresponds to strftime's %X).
	TimeFormat() string

	// AM returns the affix for ante meridiem time (corresponds to strftime's %p).
	AM() string

	// PM returns the affix for post meridiem time (corresponds to strftime's %p).
	PM() string

	// TimeAMPMFormat returns a format string to represent a time in AM/PM notation in a locale-specific way (corresponds to strftime's %r).
	TimeAMPMFormat() string

	// Day returns the name of the day of the week (corresponds to strftime's %A).
	Day(day time.Weekday) string

	// ShortDay returns the abbreviated name of the day of the week (corresponds to strftime's %a).
	ShortDay(day time.Weekday) string

	// Month returns the name of the month (corresponds to strftime's %B).
	Month(month time.Month) string

	// ShortMonth returns the abbreviated name of the month (corresponds to strftime's %b).
	ShortMonth(month time.Month) string

	// Radix returns the decimal separator between the integer and fractional parts of a decimal number.
	Radix() string

	// ThousandSeparator returns the separator character for thousands.
	ThousandSeparator() string

	// Currency returns the currency symbol and where the character should be placed relative to a monetary value.
	Currency() (symbol string, position CurrencyFormat)
}

// Default returns the current locale.
func Default() Locale {
	l, err := NewLocale("")
	if err != nil {
		panic(err)
	}
	return l
}

// CurrencyFormat describes how the currency symbol should be formatted relative to a monetary value.
type CurrencyFormat int

const (
	// BeforeValue means prefix, no separation, for example, $1.1
	BeforeValue CurrencyFormat = iota
	// AfterValue means suffix, no separation, for example, 1.1$
	AfterValue
	// BeforeValueSpace means prefix, 1-character separation, for example, $ 1.1
	BeforeValueSpace
	// AfterValueSpace means suffix, 1-character separation, for example, 1.1 $
	AfterValueSpace
	// ReplaceRadix means the currency symbol replaces the radix, for example, 1$1
	ReplaceRadix
)

func (c CurrencyFormat) String() string {
	switch c {
	case BeforeValue:
		return "BeforeValue"
	case AfterValue:
		return "AfterValue"
	case BeforeValueSpace:
		return "BeforeValueSpace"
	case AfterValueSpace:
		return "AfterValueSpace"
	case ReplaceRadix:
		return "ReplaceRadix"
	default:
		return ""
	}
}
