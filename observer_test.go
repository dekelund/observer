package observer

import (
	"errors"
	"fmt"
	"time"
)

type Client struct {
	destination string
}

type Notification struct {
	username string
}

func (event Notification) GetKey() (string, error) {
	return event.username, nil
}

func (client Client) HandleEvent(event Event) error {
	if event, ok := event.(Notification); ok {
		fmt.Println("Client:", client.destination, "received:", event.username)
	} else {
		return errors.New("Can't cast event to notification type")
	}

	return nil
}

func ExampleDistributor_ObserveAndNotify() {
	// FIXME: The test shall not rely on time.Sleep
	distributor := NewDistributor("myDistributor")
	stopObserving := distributor.ObserveAndNotify()
	defer close(stopObserving)

	distributor.RegisterObservers(
		Observe{Client{"127.0.0.1"}, []string{"bob", "alice"}},
		Observe{Client{"127.0.0.2"}, []string{"alice"}},
	)

	time.Sleep(300 * time.Microsecond)

	distributor.NotifyObservers(
		Notification{"bob"},
		Notification{"alice"},
	)

	time.Sleep(300 * time.Microsecond)

	// Output:
	// Client: 127.0.0.1 received: bob
	// Client: 127.0.0.1 received: alice
	// Client: 127.0.0.2 received: alice
}
