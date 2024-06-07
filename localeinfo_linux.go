//go:build cgo

package localeinfo

/*
#include <langinfo.h>
#include <locale.h>
#include <stdlib.h>
*/
import "C"
import (
	"runtime"
	"sync"
	"time"
	"unsafe"
)

var _ Locale = &linuxLocale{}

type linuxLocale struct {
	l      sync.Mutex
	locale C.locale_t
}

// NewLocale tries creating a new Locale from the specified locale name.
func NewLocale(name string) (Locale, error) {
	clLocale := C.CString(name)
	defer C.free(unsafe.Pointer(clLocale))
	l, err := C.newlocale(C.LC_ALL_MASK, clLocale, nil)
	if l == nil {
		return nil, err
	}
	ll := &linuxLocale{
		locale: l,
	}
	runtime.SetFinalizer(ll, func(l *linuxLocale) {
		C.freelocale(l.locale)
	})
	return ll, nil
}

func (l *linuxLocale) Encoding() string {
	l.l.Lock()
	defer l.l.Unlock()
	p := C.nl_langinfo_l(C.CODESET, l.locale)
	return C.GoString(p)
}

func (l *linuxLocale) DateTimeFormat() string {
	l.l.Lock()
	defer l.l.Unlock()
	p := C.nl_langinfo_l(C.D_T_FMT, l.locale)
	return C.GoString(p)
}

func (l *linuxLocale) DateFormat() string {
	l.l.Lock()
	defer l.l.Unlock()
	p := C.nl_langinfo_l(C.D_FMT, l.locale)
	return C.GoString(p)
}

func (l *linuxLocale) TimeFormat() string {
	l.l.Lock()
	defer l.l.Unlock()
	p := C.nl_langinfo_l(C.T_FMT, l.locale)
	return C.GoString(p)
}

func (l *linuxLocale) AM() string {
	l.l.Lock()
	defer l.l.Unlock()
	p := C.nl_langinfo_l(C.AM_STR, l.locale)
	return C.GoString(p)
}

func (l *linuxLocale) PM() string {
	l.l.Lock()
	defer l.l.Unlock()
	p := C.nl_langinfo_l(C.PM_STR, l.locale)
	return C.GoString(p)
}

func (l *linuxLocale) TimeAMPMFormat() string {
	l.l.Lock()
	defer l.l.Unlock()
	p := C.nl_langinfo_l(C.T_FMT_AMPM, l.locale)
	return C.GoString(p)
}

var days = []C.nl_item{C.DAY_1, C.DAY_2, C.DAY_3, C.DAY_4, C.DAY_5, C.DAY_6, C.DAY_7}

func (l *linuxLocale) Day(day time.Weekday) string {
	l.l.Lock()
	defer l.l.Unlock()
	p := C.nl_langinfo_l(days[day], l.locale)
	return C.GoString(p)
}

var shortDays = []C.nl_item{C.ABDAY_1, C.ABDAY_2, C.ABDAY_3, C.ABDAY_4, C.ABDAY_5, C.ABDAY_6, C.ABDAY_7}

func (l *linuxLocale) ShortDay(day time.Weekday) string {
	l.l.Lock()
	defer l.l.Unlock()
	p := C.nl_langinfo_l(shortDays[day], l.locale)
	return C.GoString(p)
}

var months = []C.nl_item{C.MON_1, C.MON_2, C.MON_3, C.MON_4, C.MON_5, C.MON_6, C.MON_7, C.MON_8, C.MON_9, C.MON_10, C.MON_11, C.MON_12}

func (l *linuxLocale) Month(month time.Month) string {
	l.l.Lock()
	defer l.l.Unlock()
	p := C.nl_langinfo_l(months[month-1], l.locale)
	return C.GoString(p)
}

var shortMonths = []C.nl_item{C.ABMON_1, C.ABMON_2, C.ABMON_3, C.ABMON_4, C.ABMON_5, C.ABMON_6, C.ABMON_7, C.ABMON_8, C.ABMON_9, C.ABMON_10, C.ABMON_11, C.ABMON_12}

func (l *linuxLocale) ShortMonth(month time.Month) string {
	l.l.Lock()
	defer l.l.Unlock()
	p := C.nl_langinfo_l(shortMonths[month-1], l.locale)
	return C.GoString(p)
}

func (l *linuxLocale) Radix() string {
	l.l.Lock()
	defer l.l.Unlock()
	p := C.nl_langinfo_l(C.RADIXCHAR, l.locale)
	return C.GoString(p)
}

func (l *linuxLocale) ThousandSeparator() string {
	l.l.Lock()
	defer l.l.Unlock()
	p := C.nl_langinfo_l(C.THOUSEP, l.locale)
	return C.GoString(p)
}

func (l *linuxLocale) Currency() (symbol string, position CurrencyFormat) {
	l.l.Lock()
	defer l.l.Unlock()
	p := C.nl_langinfo_l(C.CRNCYSTR, l.locale)
	v := C.GoString(p)
	if len(v) > 0 {
		switch v[0] {
		case '-':
			position = BeforeValue
		case '+':
			position = AfterValue
		case '.':
			position = ReplaceRadix
		}
		symbol = v[1:]
	}
	return
}
