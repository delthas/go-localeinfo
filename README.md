# go-localeinfo [![GoDoc](https://godocs.io/github.com/delthas/go-localeinfo?status.svg)](https://godocs.io/github.com/delthas/go-localeinfo)

A cross-platform library for extracting locale information.

Rather than extracting the current locale name (e.g. en_US), this library enables clients to extract monetary/numeric/time formatting information.

This library is basically a wrapper over different platform-specific calls:
- on Linux: nl_langinfo
- on Windows: GetLocaleInfoEx
- on other platforms: a stub implementation that returns empty values

## Usage

The API is well-documented in its [![GoDoc](https://godocs.io/github.com/delthas/go-localeinfo?status.svg)](https://godocs.io/github.com/delthas/go-localeinfo)

```
l := localeinfo.Default()
fmt.Printf("123%s456\n", l.ThousandSeparator())
```

## Status

Used in small-scale production environments.

- [ ] Add darwin
- [ ] Add web

## License

MIT
