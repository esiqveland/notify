package main

import (
	"fmt"
	"log"

	"github.com/esiqveland/notify"
	"github.com/godbus/dbus"
)

func main() {
	iconName := "mail-unread"

	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}

	notifier, err := notify.New(conn)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer notifier.Close()

	n := notify.Notification{
		AppName:       "Test GO App",
		ReplacesID:    uint32(0),
		AppIcon:       iconName,
		Summary:       "Test",
		Body:          "This is a test of the DBus bindings for go.",
		Actions:       []string{"cancel", "Cancel", "open", "Open"}, // tuples of (action_key, label)
		Hints:         map[string]dbus.Variant{},
		ExpireTimeout: int32(5000),
	}

	id, err := notifier.SendNotification(n)
	if err != nil {
		panic(err)
	}
	log.Printf("sent notification id: %v", id)

	actions := notifier.ActionInvoked()
	action := <-actions
	log.Printf("Action: %v Key: %v", action.Id, action.ActionKey)

	closer := <-notifier.NotificationClosed()
	log.Printf("NotificationClosed: %v Reason: %v", closer.Id, closer.Reason)

	caps, err := notifier.GetCapabilities()
	if err != nil {
		log.Printf("error fetching capabilities: %v", err)
	}
	for x := range caps {
		fmt.Printf("Registered capability: %v\n", caps[x])
	}

	info, err := notifier.GetServerInformation()
	if err != nil {
		log.Printf("error getting server information: %v", err)
	}
	fmt.Printf("Name:    %v\n", info.Name)
	fmt.Printf("Vendor:  %v\n", info.Vendor)
	fmt.Printf("Version: %v\n", info.Version)
	fmt.Printf("Spec:    %v\n", info.SpecVersion)

	// And there is a helper for just sending notifications directly if you have a connection:
	//notify.SendNotification(conn, n)

}
