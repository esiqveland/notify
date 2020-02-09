package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/esiqveland/notify"
	"github.com/godbus/dbus/v5"
)

func main() {
	wg := &sync.WaitGroup{}

	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}

	// Basic usage
	// Create a Notification to send
	iconName := "mail-unread"
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

	// Ship it!
	createdID, err := notify.SendNotification(conn, n)
	if err != nil {
		log.Printf("error sending notification: %v", err.Error())
	}
	log.Printf("created notification with id: %v", createdID)

	// List server features!
	caps, err := notify.GetCapabilities(conn)
	if err != nil {
		log.Printf("error fetching capabilities: %v", err)
	}
	for x := range caps {
		fmt.Printf("Registered capability: %v\n", caps[x])
	}

	info, err := notify.GetServerInformation(conn)
	if err != nil {
		log.Printf("error getting server information: %v", err)
	}
	fmt.Printf("Name:    %v\n", info.Name)
	fmt.Printf("Vendor:  %v\n", info.Vendor)
	fmt.Printf("Version: %v\n", info.Version)
	fmt.Printf("Spec:    %v\n", info.SpecVersion)

	// Listen for actions invoked!
	onAction := func(action *notify.ActionInvokedSignal) {
		log.Printf("ActionInvoked: %v Key: %v", action.ID, action.ActionKey)
		wg.Done()
	}

	onClosed := func(closer *notify.NotificationClosedSignal) {
		log.Printf("NotificationClosed: %v Reason: %v", closer.ID, closer.Reason)
	}

	// Notifyer interface with event delivery
	notifier, err := notify.New(conn, notify.SetOnAction(onAction), notify.SetOnClosed(onClosed), notify.SetLogger(log.New(os.Stdout, "notify: ", log.Flags())))
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer notifier.Close()

	id, err := notifier.SendNotification(n)
	if err != nil {
		log.Printf("error sending notification: %v", err)
	}
	log.Printf("sent notification id: %v", id)

	wg.Add(2)
	wg.Wait()
}
