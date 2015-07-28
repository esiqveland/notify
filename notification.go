package notify

import (
	"github.com/godbus/dbus"
	"log"
	"errors"
)

// This package is a wrapper around godbus for notification handling
// See: https://developer.gnome.org/notification-spec/
//
// Each notification displayed is allocated a unique ID by the server. (see Notify)
// This is unique within the dbus session. While the notification server is running,
// the ID will not be recycled unless the capacity of a uint32 is exceeded.
// This can be used to hide the notification before the expiration timeout is reached. (see CloseNotification)
// It can also be used to atomically replace the notification with another (Notification.ReplaceID).
// This allows you to (for instance) modify the contents of a notification while it's on-screen.
const (
	objectPath                 = "/org/freedesktop/Notifications" // the DBUS object path
	dbusNotificationsInterface = "org.freedesktop.Notifications"  // DBUS Interface
	getCapabilities            = "org.freedesktop.Notifications.GetCapabilities"
	closeNotification          = "org.freedesktop.Notifications.CloseNotification"
	notify                     = "org.freedesktop.Notifications.Notify"
	getServerInformation       = "org.freedesktop.Notifications.GetServerInformation"
)

func NewNotifier(conn *dbus.Conn) Notificator {
	return &Notifier{
		conn: conn,
	}
}

type Notifier struct {
	conn *dbus.Conn
}

// GetCapabilities gets the capabilities of the notification server.
// This call takes no parameters.
// It returns an array of strings. Each string describes an optional capability implemented by the server.
//
// Implements dbus call:
//
// STRING_ARRAY org.freedesktop.Notifications.GetCapabilities ()
//
// See also: https://developer.gnome.org/notification-spec/
func (self *Notifier) GetCapabilities() ([]string, error) {
	obj := self.conn.Object(dbusNotificationsInterface, objectPath)
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

type AskCapabilities interface {
	GetCapabilities() ([]string, error)
}

// org.freedesktop.Notifications.CloseNotification
// void org.freedesktop.Notifications.CloseNotification ( id )
//														 UINT32 id;
//
// Causes a notification to be forcefully closed and removed from the user's view. It can be used, for example, in the event that what the notification pertains to is no longer relevant, or to cancel a notification with no expiration time.
// The NotificationClosed signal is emitted by this method.
// If the notification no longer exists, an empty D-BUS Error message is sent back.
func (self *Notifier) CloseNotification(id int) (bool, error) {
	obj := self.conn.Object(dbusNotificationsInterface, objectPath)
	call := obj.Call(closeNotification, 0, uint32(id))
	if call.Err != nil {
		return false, call.Err
	}
	return true, nil
}

// Closer is an interface for implementing CloseNotification call
type Closer interface {
	CloseNotification(id int) (bool, error)
}

// ServerInformation is a holder for information returned by
// GetServerInformation call
type ServerInformation struct {
	Name        string
	Vendor      string
	Version     string
	SpecVersion string
}

// GetServerInformation
// This message returns the information on the server. Specifically, the server name, vendor, and version number.
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
func (self *Notifier) GetServerInformation() (ServerInformation, error) {
	obj := self.conn.Object(dbusNotificationsInterface, objectPath)
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

// ServerInformator is an interface for implementing GetServerInformation
type ServerInformator interface {
	GetServerInformation() (ServerInformation, error)
}

// Sends a notification to the notification server.
// Implements dbus call:
//
// UINT32 org.freedesktop.Notifications.Notify ( STRING app_name,
//												 UINT32 replaces_id,
//												 STRING app_icon,
//												 STRING summary,
//												 STRING body,
//												 ARRAY  actions,
//												 DICT   hints,
//												 INT32  expire_timeout);
//
//		Table 6. Notify Parameters
//
//		 Name	    	Type	Description
//		 app_name		STRING	The optional name of the application sending the notification. Can be blank.
//		 replaces_id	UINT32	The optional notification ID that this notification replaces. The server must atomically (ie with no flicker or other visual cues) replace the given notification with this one. This allows clients to effectively modify the notification while it's active. A value of value of 0 means that this notification won't replace any existing notifications.
//		 app_icon		STRING	The optional program icon of the calling application. See Icons and Images. Can be an empty string, indicating no icon.
//		 summary		STRING	The summary text briefly describing the notification.
//		 body			STRING	The optional detailed body text. Can be empty.
//		 actions		ARRAY	Actions are sent over as a list of pairs. Each even element in the list (starting at index 0) represents the identifier for the action. Each odd element in the list is the localized string that will be displayed to the user.
//		 hints	        DICT	Optional hints that can be passed to the server from the client program. Although clients and servers should never assume each other supports any specific hints, they can be used to pass along information, such as the process PID or window ID, that the server may be able to make use of. See Hints. Can be empty.
//       expire_timeout INT32   The timeout time in milliseconds since the display of the notification at which the notification should automatically close.
//								If -1, the notification's expiration time is dependent on the notification server's settings, and may vary for the type of notification. If 0, never expire.
//
//If replaces_id is 0, the return value is a UINT32 that represent the notification. It is unique, and will not be reused unless a MAXINT number of notifications have been generated. An acceptable implementation may just use an incrementing counter for the ID. The returned ID is always greater than zero. Servers must make sure not to return zero as an ID.
//If replaces_id is not 0, the returned value is the same value as replaces_id.
func (self *Notifier) SendNotification(n Notification) (uint32, error) {
	return SendNotification(self.conn, n)
}

func SendNotification(conn *dbus.Conn, n Notification) (uint32, error) {
	obj := conn.Object(dbusNotificationsInterface, objectPath)
	call := obj.Call(notify, 0,
		n.AppName,
		n.ReplacesID,
		n.AppIcon,
		n.Summary,
		n.Body,
		n.Actions,
		n.Hints,
		n.ExpireTimeout)
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
// Notifyer is an interface for implementing SendNotification
type Notifyer interface {
	SendNotification(n Notification) (uint32, error)
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

// Notificator is just a holder for all the small interfaces here.
type Notificator interface {
	Notifyer
	AskCapabilities
	ServerInformator
	Closer
}
