# notify

Notify is a go library for interacting with the dbus notification service defined here:
https://developer.gnome.org/notification-spec/

It can deliver notifications to desktop using dbus communication, ala how libnotify does it.
It has so far only been testing with gnome and gnome-shell 3.16 in Arch Linux. 

Please note ```notify``` is still in a very early change and no APIs are locked until a 1.0 is released.

More testers are very welcome =)

Depends on github.com/godbus/dbus.

## Quick intro

```go
package main

import (
	"github.com/esiqveland/notify"
	"github.com/godbus/dbus"
	"log"
	"fmt"
)

func main() {
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}
	
	notifier := notify.New(conn)
	
	n := notify.Notification{
		AppName:       "Test GO App",
		ReplacesID:    uint32(0),
		AppIcon:       "mail-unread",
		Summary:       "Test",
		Body:          "This is a test of the DBus bindings for go.",
		Actions:       []string{},
		Hints:         map[string]dbus.Variant{},
		ExpireTimeout: int32(5000),
	}

	id, err := notifier.SendNotification(n)
	if err != nil {
		panic(err)
	}
	log.Printf("sent notification id: %v", id)
	
	// And there is a helper for just sending notifications directly:
	notify.SendNotification(conn, n)
}

```

You should now have gotten this notification delivered to your desktop.

## TODO

- [ ] Add callback support aka dbus signals?
- [ ] Tests. I am very interested in any ideas for writing some (useful) tests for this.

## See also

`main.go` in `examples/` directory and https://developer.gnome.org/notification-spec/

## License

GPLv3
