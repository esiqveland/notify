# notify

[![GoDoc](https://godoc.org/github.com/esiqveland/notify?status.svg)](https://godoc.org/github.com/esiqveland/notify)
[![CircleCI](https://circleci.com/gh/esiqveland/notify.svg?style=svg)](https://circleci.com/gh/esiqveland/notify)

Notify is a go library for interacting with the dbus notification service defined by freedesktop.org:
https://developer.gnome.org/notification-spec/

It can deliver notifications to desktop using dbus communication, ala how libnotify does it.
It has so far only been tested on gnome and gnome-shell 3.16/3.18/3.22 in Arch Linux. 

Please note ```notify``` is still in a very early change and no APIs are locked until a 1.0 is released.

More testers are very welcome =)

Depends on:
 - [godbus](https://github.com/godbus/dbus).

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
