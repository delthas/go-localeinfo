# go-localeinfo [![GoDoc](https://godocs.io/github.com/delthas/go-localeinfo?status.svg)](https://godocs.io/github.com/delthas/go-localeinfo) [![stability-experimental](https://img.shields.io/badge/stability-experimental-orange.svg)](https://github.com/emersion/stability-badges#experimental)

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

The API could be slightly changed in backwards-incompatible ways for now.

- [ ] Add darwin
- [ ] Add web

## License

MIT
