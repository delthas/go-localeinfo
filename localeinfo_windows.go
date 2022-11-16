//go:build cgo

package localeinfo

/*
#define WIN32_LEAN_AND_MEAN
#include <Windows.h>
*/
import "C"
import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
	"unicode/utf16"
	"unsafe"
)

var _ Locale = &windowsLocale{}

type windowsLocale struct {
	l      sync.Mutex
	buf    *uint16
	locale *C.ushort
}

const bufSize = 1024

func NewLocale(name string) (*windowsLocale, error) {
	var cLocale *uint16
	if name == "" {
		// LOCALE_NAME_USER_DEFAULT is NULL
		cLocale = nil
	} else {
		nameUTF16 := utf16.Encode([]rune(name + "\x00"))
		cLocale = &nameUTF16[0]
	}
	buf := make([]uint16, bufSize)
	return &windowsLocale{
		buf:    &buf[0],
		locale: (*C.ushort)(unsafe.Pointer(cLocale)),
	}, nil
}

func (l *windowsLocale) decode(n int) string {
	if n == 0 {
		return ""
	}
	var data []uint16
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	sh.Data = uintptr(unsafe.Pointer(l.buf))
	sh.Len = n - 1
	sh.Cap = n - 1
	s := string(utf16.Decode(data))
	runtime.KeepAlive(l.buf)
	return s
}

func (l *windowsLocale) localeInfo(t C.ulong) string {
	l.l.Lock()
	defer l.l.Unlock()
	n, err := C.GetLocaleInfoEx(l.locale, t, (*C.ushort)(unsafe.Pointer(l.buf)), bufSize)
	if err != nil {
		return ""
	}
	return l.decode(int(n))
}

func (l *windowsLocale) localeInfoInt(t C.ulong, def int) int {
	var n int32
	_, err := C.GetLocaleInfoEx(l.locale, t|C.LOCALE_RETURN_NUMBER, (*C.ushort)(unsafe.Pointer(&n)), 2)
	if err != nil {
		return def
	}
	return int(n)
}

func convertDateTimeFormat(fmt string) string {
	var sb strings.Builder
	sb.Grow(len(fmt))
	var current = '\x00'
	var count = 0
	escape := -1
	for i, c := range fmt + "\x00" {
		if escape >= 0 && c != '\x00' {
			if c == '\'' {
				if escape == i-1 {
					escape = -1
					sb.WriteRune(c)
					continue
				}
				escape = -1
				continue
			}
			if c == '%' {
				sb.WriteRune('%')
			}
			sb.WriteRune(c)
			continue
		}
		if count > 0 && c == current {
			count++
			continue
		}
		if count > 0 {
			if current == 'h' {
				sb.WriteString("%I")
			} else if current == 'H' {
				sb.WriteString("%H")
			} else if current == 'm' {
				sb.WriteString("%M")
			} else if current == 's' {
				sb.WriteString("%S")
			} else if current == 't' {
				sb.WriteString("%p")
			} else if current == 'd' && count == 1 {
				// zero-padding is better than space-padding when we want no padding
				// otherwise we would get eg 01/ 1
				sb.WriteString("%d")
			} else if current == 'd' && count == 2 {
				sb.WriteString("%d")
			} else if current == 'd' && count == 3 {
				sb.WriteString("%a")
			} else if current == 'd' && count == 4 {
				sb.WriteString("%A")
			} else if current == 'M' && count <= 2 {
				sb.WriteString("%m")
			} else if current == 'M' && count == 3 {
				sb.WriteString("%b")
			} else if current == 'M' && count == 4 {
				sb.WriteString("%B")
			} else if current == 'y' && count <= 2 {
				sb.WriteString("%y")
			} else if current == 'y' && count >= 4 {
				sb.WriteString("%Y")
			}
			count = 0
		}
		switch c {
		case '\x00':
			return sb.String()
		case '\'':
			escape = i
		case 'd', 'M', 'y', 'h', 'H', 'm', 's', 't':
			current = c
			count = 1
		case '%':
			sb.WriteString("%%")
		default:
			sb.WriteRune(c)
		}
	}
	panic("unreachable")
}

func (l *windowsLocale) Encoding() string {
	return l.localeInfo(C.LOCALE_SNAME)
}

func (l *windowsLocale) DateTimeFormat() string {
	// Best-effort, because there is no equivalent of %c.
	// ReactOS outputs date and time separated by a space, let's do that too.
	d := l.DateFormat()
	t := l.TimeFormat()
	if d == "" || t == "" {
		return ""
	}
	return fmt.Sprintf("%s %s", d, t)
}

func (l *windowsLocale) DateFormat() string {
	// https://learn.microsoft.com/en-us/windows/win32/intl/day--month--year--and-era-format-pictures
	return convertDateTimeFormat(l.localeInfo(C.LOCALE_SSHORTDATE))
}

func (l *windowsLocale) TimeFormat() string {
	// https://learn.microsoft.com/en-us/windows/win32/intl/hour--minute--and-second-format-pictures
	return convertDateTimeFormat(l.localeInfo(C.LOCALE_STIMEFORMAT))
}

func (l *windowsLocale) AM() string {
	return l.localeInfo(C.LOCALE_S1159)
}

func (l *windowsLocale) PM() string {
	return l.localeInfo(C.LOCALE_S2359)
}

func (l *windowsLocale) TimeAMPMFormat() string {
	// Best-effort: use TimeFormat
	return l.TimeFormat()
}

var days = []C.ulong{C.LOCALE_SDAYNAME7, C.LOCALE_SDAYNAME1, C.LOCALE_SDAYNAME2, C.LOCALE_SDAYNAME3, C.LOCALE_SDAYNAME4, C.LOCALE_SDAYNAME5, C.LOCALE_SDAYNAME6}

func (l *windowsLocale) Day(day time.Weekday) string {
	return l.localeInfo(days[day])
}

var shortDays = []C.ulong{C.LOCALE_SABBREVDAYNAME7, C.LOCALE_SABBREVDAYNAME1, C.LOCALE_SABBREVDAYNAME2, C.LOCALE_SABBREVDAYNAME3, C.LOCALE_SABBREVDAYNAME4, C.LOCALE_SABBREVDAYNAME5, C.LOCALE_SABBREVDAYNAME6}

func (l *windowsLocale) ShortDay(day time.Weekday) string {
	return l.localeInfo(shortDays[day])
}

var months = []C.ulong{C.LOCALE_SMONTHNAME1, C.LOCALE_SMONTHNAME2, C.LOCALE_SMONTHNAME3, C.LOCALE_SMONTHNAME4, C.LOCALE_SMONTHNAME5, C.LOCALE_SMONTHNAME6, C.LOCALE_SMONTHNAME7, C.LOCALE_SMONTHNAME8, C.LOCALE_SMONTHNAME9, C.LOCALE_SMONTHNAME10, C.LOCALE_SMONTHNAME11, C.LOCALE_SMONTHNAME12}

func (l *windowsLocale) Month(month time.Month) string {
	return l.localeInfo(months[month-1])
}

var shortMonths = []C.ulong{C.LOCALE_SABBREVMONTHNAME1, C.LOCALE_SABBREVMONTHNAME2, C.LOCALE_SABBREVMONTHNAME3, C.LOCALE_SABBREVMONTHNAME4, C.LOCALE_SABBREVMONTHNAME5, C.LOCALE_SABBREVMONTHNAME6, C.LOCALE_SABBREVMONTHNAME7, C.LOCALE_SABBREVMONTHNAME8, C.LOCALE_SABBREVMONTHNAME9, C.LOCALE_SABBREVMONTHNAME10, C.LOCALE_SABBREVMONTHNAME11, C.LOCALE_SABBREVMONTHNAME12}

func (l *windowsLocale) ShortMonth(month time.Month) string {
	return l.localeInfo(shortMonths[month-1])
}

func (l *windowsLocale) Radix() string {
	return l.localeInfo(C.LOCALE_SDECIMAL)
}

func (l *windowsLocale) ThousandSeparator() string {
	return l.localeInfo(C.LOCALE_STHOUSAND)
}

func (l *windowsLocale) Currency() (symbol string, position CurrencyFormat) {
	symbol = l.localeInfo(C.LOCALE_SCURRENCY)
	switch l.localeInfoInt(C.LOCALE_ICURRENCY, 0) {
	case 0:
		position = BeforeValue
	case 1:
		position = AfterValue
	case 2:
		position = BeforeValueSpace
	case 3:
		position = AfterValueSpace
	}
	return
}
