package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/godbus/dbus/v5"

	"github.com/esiqveland/notify"
)

func main() {
	err := runMain()
	if err != nil {
		log.Printf("\nerror: %v\n", err)
		os.Exit(1)
	}
}

func runMain() error {
	wg := &sync.WaitGroup{}

	conn, err := dbus.SessionBusPrivate()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	if err = conn.Auth(nil); err != nil {
		panic(err)
	}

	if err = conn.Hello(); err != nil {
		panic(err)
	}

	DebugServerFeatures(conn)

	// Basic usage
	soundHint := notify.HintSoundWithName(
		//"message-new-instant",
		"trash-empty",
	)

	// Create a Notification to send
	iconName := "mail-unread"
	n := notify.Notification{
		AppName:    "Test GO App",
		ReplacesID: uint32(0),
		AppIcon:    iconName,
		Summary:    "Test",
		Body:       "This is a test of the DBus bindings for go with sound.",
		Actions: []notify.Action{
			{Key: "cancel", Label: "Cancel"},
			{Key: "open", Label: "Open"},
		},
		Hints: map[string]dbus.Variant{
			soundHint.ID: soundHint.Variant,
		},
		ExpireTimeout: time.Second * 5,
	}
	n.SetUrgency(notify.UrgencyCritical)

	counter := int32(0)
	// Listen for actions invoked!
	onAction := func(action *notify.ActionInvokedSignal) {
		atomic.AddInt32(&counter, 1)
		log.Printf("ActionInvoked: %v Key: %v", action.ID, action.ActionKey)
		wg.Done()
	}

	onClosed := func(closer *notify.NotificationClosedSignal) {
		atomic.AddInt32(&counter, 1)
		log.Printf("NotificationClosed: %v Reason: %v", closer.ID, closer.Reason)
		wg.Done()
	}

	// Notifier instance with event delivery:
	notifier, err := notify.New(
		conn,
		// action event handler
		notify.WithOnAction(onAction),
		// closed event handler
		notify.WithOnClosed(onClosed),
		// override with custom logger
		notify.WithLogger(log.New(os.Stdout, "notify: ", log.Flags())),
	)
	if err != nil {
		return err
	}
	defer notifier.Close()

	rgbaSample, err := readImage("./big.jpg")
	if err != nil {
		return err
	}

	absFilePath, err := filepath.Abs("./small.png")
	if err != nil {
		return err
	}
	hintProfileImage := notify.HintImageDataRGBA(rgbaSample)
	hintImageByPath := notify.HintImageFilePath(absFilePath)
	urgencyHint := notify.HintUrgency(notify.UrgencyCritical)
	n.AddHint(urgencyHint)
	n.AddHint(hintImageByPath)
	// according to spec, image-data hint should have precedence:
	n.AddHint(hintProfileImage)

	wg.Add(1)
	id, err := notifier.SendNotification(n)
	if err != nil {
		log.Printf("error sending notification: %v", err)
	}
	log.Printf("sent notification id: %v", id)

	// add two extra so the process should block forever:
	// this is due to a bug in Gnome that delivers multiple copies of the action signal...
	wg.Add(2)
	wg.Wait()

	log.Printf("total signal count received: %d", atomic.LoadInt32(&counter))

	return nil
}

func DebugServerFeatures(conn *dbus.Conn) {
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
}

func readImage(path string) (*image.RGBA, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	decode, _, err := image.Decode(fd)
	if err != nil {
		return nil, err
	}
	img, ok := decode.(*image.RGBA)
	if !ok {
		b := decode.Bounds()
		m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(m, m.Bounds(), decode, b.Min, draw.Src)
		return m, nil
		//return nil, fmt.Errorf("decode RGBA failed of: %v", path)
	}
	return img, err
}
