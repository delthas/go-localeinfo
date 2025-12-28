package localeinfo

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode/utf16"
	"unsafe"
)

var (
	kernel32        = syscall.NewLazyDLL("kernel32.dll")
	getLocaleInfoEx = kernel32.NewProc("GetLocaleInfoEx")
)

const (
	cRETURN_NUMBER      = 0x20000000
	cSNAME              = 0x0000005c
	cSSHORTDATE         = 0x0000001F
	cSTIMEFORMAT        = 0x00001003
	cSAM                = 0x00000028
	cSPM                = 0x00000029
	cSDAYNAME1          = 0x0000002A
	cSDAYNAME2          = 0x0000002B
	cSDAYNAME3          = 0x0000002C
	cSDAYNAME4          = 0x0000002D
	cSDAYNAME5          = 0x0000002E
	cSDAYNAME6          = 0x0000002F
	cSDAYNAME7          = 0x00000030
	cSABBREVDAYNAME1    = 0x00000031
	cSABBREVDAYNAME2    = 0x00000032
	cSABBREVDAYNAME3    = 0x00000033
	cSABBREVDAYNAME4    = 0x00000034
	cSABBREVDAYNAME5    = 0x00000035
	cSABBREVDAYNAME6    = 0x00000036
	cSABBREVDAYNAME7    = 0x00000037
	cSMONTHNAME1        = 0x00000038
	cSMONTHNAME2        = 0x00000039
	cSMONTHNAME3        = 0x0000003A
	cSMONTHNAME4        = 0x0000003B
	cSMONTHNAME5        = 0x0000003C
	cSMONTHNAME6        = 0x0000003D
	cSMONTHNAME7        = 0x0000003E
	cSMONTHNAME8        = 0x0000003F
	cSMONTHNAME9        = 0x00000040
	cSMONTHNAME10       = 0x00000041
	cSMONTHNAME11       = 0x00000042
	cSMONTHNAME12       = 0x00000043
	cSABBREVMONTHNAME1  = 0x00000044
	cSABBREVMONTHNAME2  = 0x00000045
	cSABBREVMONTHNAME3  = 0x00000046
	cSABBREVMONTHNAME4  = 0x00000047
	cSABBREVMONTHNAME5  = 0x00000048
	cSABBREVMONTHNAME6  = 0x00000049
	cSABBREVMONTHNAME7  = 0x0000004A
	cSABBREVMONTHNAME8  = 0x0000004B
	cSABBREVMONTHNAME9  = 0x0000004C
	cSABBREVMONTHNAME10 = 0x0000004D
	cSABBREVMONTHNAME11 = 0x0000004E
	cSABBREVMONTHNAME12 = 0x0000004F
	cSDECIMAL           = 0x0000000E
	cSTHOUSAND          = 0x0000000F
	cSCURRENCY          = 0x00000014
	cICURRENCY          = 0x0000001B
)

var _ Locale = &windowsLocale{}

type windowsLocale struct {
	l      sync.Mutex
	buf    *uint16
	locale *uint16
}

const bufSize = 1024

func NewLocale(name string) (Locale, error) {
	var cLocale *uint16
	if name == "" {
		// cNAME_USER_DEFAULT is NULL
		cLocale = nil
	} else {
		nameUTF16 := utf16.Encode([]rune(name + "\x00"))
		cLocale = &nameUTF16[0]
	}
	buf := make([]uint16, bufSize)
	return &windowsLocale{
		buf:    &buf[0],
		locale: cLocale,
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

func (l *windowsLocale) localeInfo(t uint32) string {
	l.l.Lock()
	defer l.l.Unlock()
	r, _, _ := getLocaleInfoEx.Call(uintptr(unsafe.Pointer(l.locale)), uintptr(t), uintptr(unsafe.Pointer(l.buf)), uintptr(bufSize))
	if r <= 0 {
		return ""
	}
	return l.decode(int(r))
}

func (l *windowsLocale) localeInfoInt(t uint32, def int) int {
	var n int32
	r, _, _ := getLocaleInfoEx.Call(uintptr(unsafe.Pointer(l.locale)), uintptr(t|cRETURN_NUMBER), uintptr(unsafe.Pointer(&n)), 2)
	if r <= 0 {
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
	return l.localeInfo(cSNAME)
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
	return convertDateTimeFormat(l.localeInfo(cSSHORTDATE))
}

func (l *windowsLocale) TimeFormat() string {
	// https://learn.microsoft.com/en-us/windows/win32/intl/hour--minute--and-second-format-pictures
	return convertDateTimeFormat(l.localeInfo(cSTIMEFORMAT))
}

func (l *windowsLocale) AM() string {
	return l.localeInfo(cSAM)
}

func (l *windowsLocale) PM() string {
	return l.localeInfo(cSPM)
}

func (l *windowsLocale) TimeAMPMFormat() string {
	// Best-effort: use TimeFormat
	return l.TimeFormat()
}

var days = []uint32{cSDAYNAME7, cSDAYNAME1, cSDAYNAME2, cSDAYNAME3, cSDAYNAME4, cSDAYNAME5, cSDAYNAME6}

func (l *windowsLocale) Day(day time.Weekday) string {
	return l.localeInfo(days[day])
}

var shortDays = []uint32{cSABBREVDAYNAME7, cSABBREVDAYNAME1, cSABBREVDAYNAME2, cSABBREVDAYNAME3, cSABBREVDAYNAME4, cSABBREVDAYNAME5, cSABBREVDAYNAME6}

func (l *windowsLocale) ShortDay(day time.Weekday) string {
	return l.localeInfo(shortDays[day])
}

var months = []uint32{cSMONTHNAME1, cSMONTHNAME2, cSMONTHNAME3, cSMONTHNAME4, cSMONTHNAME5, cSMONTHNAME6, cSMONTHNAME7, cSMONTHNAME8, cSMONTHNAME9, cSMONTHNAME10, cSMONTHNAME11, cSMONTHNAME12}

func (l *windowsLocale) Month(month time.Month) string {
	return l.localeInfo(months[month-1])
}

var shortMonths = []uint32{cSABBREVMONTHNAME1, cSABBREVMONTHNAME2, cSABBREVMONTHNAME3, cSABBREVMONTHNAME4, cSABBREVMONTHNAME5, cSABBREVMONTHNAME6, cSABBREVMONTHNAME7, cSABBREVMONTHNAME8, cSABBREVMONTHNAME9, cSABBREVMONTHNAME10, cSABBREVMONTHNAME11, cSABBREVMONTHNAME12}

func (l *windowsLocale) ShortMonth(month time.Month) string {
	return l.localeInfo(shortMonths[month-1])
}

func (l *windowsLocale) Radix() string {
	return l.localeInfo(cSDECIMAL)
}

func (l *windowsLocale) ThousandSeparator() string {
	return l.localeInfo(cSTHOUSAND)
}

func (l *windowsLocale) Currency() (symbol string, position CurrencyFormat) {
	symbol = l.localeInfo(cSCURRENCY)
	switch l.localeInfoInt(cICURRENCY, 0) {
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
