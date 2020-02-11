/*
Package notify is a wrapper around godbus for dbus notification interface
See: https://developer.gnome.org/notification-spec/ and
https://github.com/godbus/dbus

The package provides exported methods for simple usage, e.g. just show a notification.
It also provides the interface Notifier that includes signal delivery for notifications on the system.

Each notification created is allocated a unique ID by the server (see Notification).
This ID is unique within the dbus session, and is an increasing number.
While the notification server is running, the ID will not be recycled unless the capacity of a uint32 is exceeded.

The ID can be used to hide the notification before the expiration timeout is reached (see CloseNotification).

The ID can also be used to atomically replace the notification with another (Notification.ReplaceID).
This allows you to (for instance) modify the contents of a notification while it's on-screen.
*/
package notify
