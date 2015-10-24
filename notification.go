package notify

import (
	"errors"
	"log"

	"github.com/godbus/dbus"
)

const (
	objectPath                    = "/org/freedesktop/Notifications" // the DBUS object path
	dbusNotificationsInterface    = "org.freedesktop.Notifications"  // DBUS Interface
	dbusNotificationClosed        = "org.freedesktop.Notifications.NotificationClosed"
	dbusNotificationActionInvoked = "org.freedesktop.Notifications.ActionInvoked"
	getCapabilities               = "org.freedesktop.Notifications.GetCapabilities"
	closeNotification             = "org.freedesktop.Notifications.CloseNotification"
	notify                        = "org.freedesktop.Notifications.Notify"
	getServerInformation          = "org.freedesktop.Notifications.GetServerInformation"
)

// New creates a new Notificator using conn
// New sets up a Notifier that listenes on dbus' signals regarding
// Notifications, e.g.
func New(conn *dbus.Conn) Notifier {
	n := &notifier{
		conn:   conn,
		signal: make(chan *dbus.Signal, 5),
		closer: make(chan *NotificationClosedSignal, 5),
		action: make(chan *ActionInvokedSignal, 5),
		done:   make(chan bool),
	}

	n.conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
		"type='signal',path='"+objectPath+"',interface='"+dbusNotificationsInterface+"'")

	go (func() {
		received := 0
		for {
			select {
			case signal := <-n.signal:
				received += 1
				log.Printf("got signal: %v Signal: %+v", received, signal)
				n.handleSignal(signal)
			// its all over, exit and go home
			case <-n.done:
				log.Printf("its all over, go home")
				return
			}
		}
	})()

	// setup signal reception
	n.conn.Signal(n.signal)

	return n
}

func (n notifier) handleSignal(signal *dbus.Signal) {
	switch signal.Name {
	case dbusNotificationClosed:
		n.closer <- &NotificationClosedSignal{
			Id:     signal.Body[0].(uint32),
			Reason: signal.Body[1].(uint32),
		}
	case dbusNotificationActionInvoked:
		n.action <- &ActionInvokedSignal{
			Id:        signal.Body[0].(uint32),
			ActionKey: signal.Body[1].(string),
		}
	default:
		log.Printf("unknown signal: %+v", signal)
	}
}

// notifier implements Notificator
type notifier struct {
	conn   *dbus.Conn
	signal chan *dbus.Signal
	closer chan *NotificationClosedSignal
	action chan *ActionInvokedSignal
	done   chan bool
}

// Notification holds all information needed for creating a notification
type Notification struct {
	AppName       string
	ReplacesID    uint32
	AppIcon       string // see icons here: http://standards.freedesktop.org/icon-naming-spec/icon-naming-spec-latest.html
	Summary       string
	Body          string
	Actions       []string
	Hints         map[string]dbus.Variant
	ExpireTimeout int32 // millisecond to show notification
}

// SendNotification sends a notification to the notification server.
// Implements dbus call:
//
//     UINT32 org.freedesktop.Notifications.Notify ( STRING app_name,
//	    										 UINT32 replaces_id,
//	    										 STRING app_icon,
//	    										 STRING summary,
//	    										 STRING body,
//	    										 ARRAY  actions,
//	    										 DICT   hints,
//	    										 INT32  expire_timeout);
//
//		Name	    	Type	Description
//		app_name		STRING	The optional name of the application sending the notification. Can be blank.
//		replaces_id	    UINT32	The optional notification ID that this notification replaces. The server must atomically (ie with no flicker or other visual cues) replace the given notification with this one. This allows clients to effectively modify the notification while it's active. A value of value of 0 means that this notification won't replace any existing notifications.
//		app_icon		STRING	The optional program icon of the calling application. Can be an empty string, indicating no icon.
//		summary		    STRING	The summary text briefly describing the notification.
//		body			STRING	The optional detailed body text. Can be empty.
//		actions		    ARRAY	Actions are sent over as a list of pairs. Each even element in the list (starting at index 0) represents the identifier for the action. Each odd element in the list is the localized string that will be displayed to the user.
//		hints	        DICT	Optional hints that can be passed to the server from the client program. Although clients and servers should never assume each other supports any specific hints, they can be used to pass along information, such as the process PID or window ID, that the server may be able to make use of. See Hints. Can be empty.
//      expire_timeout  INT32   The timeout time in milliseconds since the display of the notification at which the notification should automatically close.
//								If -1, the notification's expiration time is dependent on the notification server's settings, and may vary for the type of notification. If 0, never expire.
//
// If replaces_id is 0, the return value is a UINT32 that represent the notification. It is unique, and will not be reused unless a MAXINT number of notifications have been generated. An acceptable implementation may just use an incrementing counter for the ID. The returned ID is always greater than zero. Servers must make sure not to return zero as an ID.
// If replaces_id is not 0, the returned value is the same value as replaces_id.
func (n *notifier) SendNotification(note Notification) (uint32, error) {
	return SendNotification(n.conn, note)
}

// SendNotification is same as Notifier.SendNotification
// Provided for convenience.
func SendNotification(conn *dbus.Conn, note Notification) (uint32, error) {
	obj := conn.Object(dbusNotificationsInterface, objectPath)
	call := obj.Call(notify, 0,
		note.AppName,
		note.ReplacesID,
		note.AppIcon,
		note.Summary,
		note.Body,
		note.Actions,
		note.Hints,
		note.ExpireTimeout)
	if call.Err != nil {
		return 0, call.Err
	}
	var ret uint32
	err := call.Store(&ret)
	if err != nil {
		log.Printf("error getting uint32 ret value: %v", err)
		return ret, err
	}
	return ret, nil
}

// GetCapabilities gets the capabilities of the notification server.
// This call takes no parameters.
// It returns an array of strings. Each string describes an optional capability implemented by the server.
//
// See also: https://developer.gnome.org/notification-spec/
func (n *notifier) GetCapabilities() ([]string, error) {
	obj := n.conn.Object(dbusNotificationsInterface, objectPath)
	call := obj.Call(getCapabilities, 0)
	if call.Err != nil {
		log.Printf("error calling GetCapabilities: %v", call.Err)
		return []string{}, call.Err
	}
	var ret []string
	err := call.Store(&ret)
	if err != nil {
		log.Printf("error getting capabilities ret value: %v", err)
		return ret, err
	}
	return ret, nil
}

// CloseNotification causes a notification to be forcefully closed and removed from the user's view.
// It can be used, for example, in the event that what the notification pertains to is no longer relevant,
// or to cancel a notification with no expiration time.
//
// The NotificationClosed (dbus) signal is emitted by this method.
// If the notification no longer exists, an empty D-BUS Error message is sent back.
func (n *notifier) CloseNotification(id int) (bool, error) {
	obj := n.conn.Object(dbusNotificationsInterface, objectPath)
	call := obj.Call(closeNotification, 0, uint32(id))
	if call.Err != nil {
		return false, call.Err
	}
	return true, nil
}

// ServerInformation is a holder for information returned by
// GetServerInformation call.
type ServerInformation struct {
	Name        string
	Vendor      string
	Version     string
	SpecVersion string
}

// GetServerInformation returns the information on the server.
//
// org.freedesktop.Notifications.GetServerInformation
//
//  GetServerInformation Return Values
//
//		Name		 Type	  Description
//		name		 STRING	  The product name of the server.
//		vendor		 STRING	  The vendor name. For example, "KDE," "GNOME," "freedesktop.org," or "Microsoft."
//		version		 STRING	  The server's version number.
//		spec_version STRING	  The specification version the server is compliant with.
//
func (n *notifier) GetServerInformation() (ServerInformation, error) {
	obj := n.conn.Object(dbusNotificationsInterface, objectPath)
	if obj == nil {
		return ServerInformation{}, errors.New("error creating dbus call object")
	}
	call := obj.Call(getServerInformation, 0)
	if call.Err != nil {
		log.Printf("Error calling %v: %v", getServerInformation, call.Err)
		return ServerInformation{}, call.Err
	}

	ret := ServerInformation{}
	err := call.Store(&ret.Name, &ret.Vendor, &ret.Version, &ret.SpecVersion)
	if err != nil {
		log.Printf("error reading %v return values: %v", getServerInformation, err)
		return ret, err
	}
	return ret, nil
}

type NotificationClosedSignal struct {
	Id     uint32
	Reason uint32
}

func (n *notifier) NotificationClosed() <-chan *NotificationClosedSignal {
	return n.closer
}

type ActionInvokedSignal struct {
	Id        uint32
	ActionKey string
}

func (n *notifier) ActionInvoked() <-chan *ActionInvokedSignal {
	return n.action
}

// Close cleans up and shuts down signal delivery loop
func (n *notifier) Close() error {
	// remove signal reception
	n.done <- true
	defer n.conn.Signal(n.signal)
	close(n.closer)
	close(n.action)
	close(n.done)
	err := n.conn.Close()
	return err
}

// Notifier is an interface for implementing the operations supported by the
// freedesktop DBus Notifications object.
type Notifier interface {
	SendNotification(n Notification) (uint32, error)
	GetCapabilities() ([]string, error)
	GetServerInformation() (ServerInformation, error)
	CloseNotification(id int) (bool, error)
	NotificationClosed() <-chan *NotificationClosedSignal
	ActionInvoked() <-chan *ActionInvokedSignal
	Close() error
}
