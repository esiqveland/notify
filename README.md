# notify

[![GoDoc](https://godoc.org/github.com/esiqveland/notify?status.svg)](https://godoc.org/github.com/esiqveland/notify)
[![CircleCI](https://circleci.com/gh/esiqveland/notify.svg?style=svg)](https://circleci.com/gh/esiqveland/notify)

Notify is a go library for interacting with the dbus notification service defined by freedesktop.org:
https://developer.gnome.org/notification-spec/

`notify` can deliver desktop notifications over dbus, ala how libnotify does it.

Please note `notify` is still in motion and APIs are not locked until a 1.0 is released.

More testers are very welcome =)

Depends on:
 - [godbus](https://github.com/godbus/dbus).

## Changelog

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

## License

GPLv3
