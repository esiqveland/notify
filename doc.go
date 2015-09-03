/*
The notify package is a wrapper around godbus for dbus notification interface
See: https://developer.gnome.org/notification-spec/ and
https://github.com/godbus/dbus

Each notification displayed is allocated a unique ID by the server. (see Notify)
This ID unique within the dbus session. While the notification server is running,
the ID will not be recycled unless the capacity of a uint32 is exceeded.

This can be used to hide the notification before the expiration timeout is reached. (see CloseNotification)

The ID can also be used to atomically replace the notification with another (Notification.ReplaceID).
This allows you to (for instance) modify the contents of a notification while it's on-screen.
*/
package notify
