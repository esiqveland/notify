
<div align="right">
  <details>
    <summary >🌐 Language</summary>
    <div>
      <div align="center">
        <a href="https://openaitx.github.io/view.html?user=esiqveland&project=notify&lang=en">English</a>
        | <a href="https://openaitx.github.io/view.html?user=esiqveland&project=notify&lang=zh-CN">简体中文</a>
        | <a href="https://openaitx.github.io/view.html?user=esiqveland&project=notify&lang=zh-TW">繁體中文</a>
        | <a href="https://openaitx.github.io/view.html?user=esiqveland&project=notify&lang=ja">日本語</a>
        | <a href="https://openaitx.github.io/view.html?user=esiqveland&project=notify&lang=ko">한국어</a>
        | <a href="https://openaitx.github.io/view.html?user=esiqveland&project=notify&lang=hi">हिन्दी</a>
        | <a href="https://openaitx.github.io/view.html?user=esiqveland&project=notify&lang=th">ไทย</a>
        | <a href="https://openaitx.github.io/view.html?user=esiqveland&project=notify&lang=fr">Français</a>
        | <a href="https://openaitx.github.io/view.html?user=esiqveland&project=notify&lang=de">Deutsch</a>
        | <a href="https://openaitx.github.io/view.html?user=esiqveland&project=notify&lang=es">Español</a>
        | <a href="https://openaitx.github.io/view.html?user=esiqveland&project=notify&lang=it">Italiano</a>
        | <a href="https://openaitx.github.io/view.html?user=esiqveland&project=notify&lang=ru">Русский</a>
        | <a href="https://openaitx.github.io/view.html?user=esiqveland&project=notify&lang=pt">Português</a>
        | <a href="https://openaitx.github.io/view.html?user=esiqveland&project=notify&lang=nl">Nederlands</a>
        | <a href="https://openaitx.github.io/view.html?user=esiqveland&project=notify&lang=pl">Polski</a>
        | <a href="https://openaitx.github.io/view.html?user=esiqveland&project=notify&lang=ar">العربية</a>
        | <a href="https://openaitx.github.io/view.html?user=esiqveland&project=notify&lang=fa">فارسی</a>
        | <a href="https://openaitx.github.io/view.html?user=esiqveland&project=notify&lang=tr">Türkçe</a>
        | <a href="https://openaitx.github.io/view.html?user=esiqveland&project=notify&lang=vi">Tiếng Việt</a>
        | <a href="https://openaitx.github.io/view.html?user=esiqveland&project=notify&lang=id">Bahasa Indonesia</a>
        | <a href="https://openaitx.github.io/view.html?user=esiqveland&project=notify&lang=as">অসমীয়া</
      </div>
    </div>
  </details>
</div>

# notify

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/esiqveland/notify?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/esiqveland/notify)](https://goreportcard.com/report/github.com/esiqveland/notify)
![Build](https://github.com/esiqveland/notify/actions/workflows/go.yml/badge.svg?branch=master)

Notify is a go library for interacting with the dbus notification service defined by freedesktop.org:
https://specifications.freedesktop.org/notification/latest/

`notify` can deliver desktop notifications over dbus, ala how libnotify does it.

Please note `notify` is still in motion and APIs are not locked until a 1.0 is released.

More testers are very welcome =)

Depends on:
 - [godbus](https://github.com/godbus/dbus).

## Changelog
- v0.11.2: Introduce helpers ExpireTimeoutSetByNotificationServer, ExpireTimeoutNever
- v0.11.1: Fix a race during Close() #11
- v0.11.0: re-release under BSD license
- v0.10.0: stricter types: [some breaking changes](https://github.com/esiqveland/notify/releases/tag/v0.10.0)
- v0.9.0: [some breaking changes](https://github.com/esiqveland/notify/releases/tag/v0.9.0)
- v0.2.1: dbus: gomod: lock to dbus v5
- v0.2.0: `Notifier.Close()` no longer calls `.Close()` on the underlying `dbus.Conn`

## Quick intro
See example: [main.go](https://github.com/esiqveland/notify/blob/master/example/main.go).

Clone repo and go to examples folder:

``` go run main.go ```


## TODO

- [x] Add callback support aka dbus signals.
- [ ] Tests. I am very interested in any ideas for writing some (useful) tests for this.

## See also

The Gnome notification spec https://developer.gnome.org/notification-spec/.


## Contributors
Thanks to user [emersion](https://github.com/emersion) for great ideas on receiving signals.

Thanks to [Merovius](https://github.com/Merovius) for fixing race during Close().

## License

BSD 3-Clause

