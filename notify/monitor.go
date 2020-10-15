package notify

import "fmt"

// Monitor is a simple wrapper to send a notification to a Notifier service in the event of a panic.  It then
// raises the panic message after flushing the Notification.
//
// Example
//     defer notify.Monitor(myNotifier)
func Monitor(n Notifier) {
	if err := recover(); err != nil {
		if n != nil {
			e, ok := err.(error)
			if !ok {
				e = fmt.Errorf("%v", err) //nolint
			}

			_, _ = n.Send(e)
			n.Flush()
		}

		panic(err)
	}
}
