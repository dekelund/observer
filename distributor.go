package observer

type dist struct {
	name         string
	observers    map[string][]Observer
	events       chan Event
	toRegister   chan Observe
	toUnregister chan Observe
}

// Distributor provides means to register and notify events
// to instances that implements the interface Observer.
type Distributor struct {
	*dist
}

// NewDistributor instantiates a Distributor
func NewDistributor(name string) Distributor {
	observers := map[string][]Observer{}

	return Distributor{&dist{name, observers, nil, nil, nil}}
}

// ObserveAndNotify sets up a go-routine mapping incoming
// events to registered observers.
// It returns a channel, use the channel to stop observing,
// either by writing a boolean to the channel, or by closing
// the channel.
func (distributor Distributor) ObserveAndNotify() chan bool {
	stopObserving := make(chan bool, 1)

	distributor.events = make(chan Event, 1024)
	distributor.toRegister = make(chan Observe, 8)
	distributor.toUnregister = make(chan Observe, 8)

	go func() {
		defer close(distributor.events)
		defer close(distributor.toRegister)
		defer close(distributor.toUnregister)

		var event Event
		var observer Observe

		ok := true
		for ok {
			select {
			case event, ok = <-distributor.events:
				distributor.notifyObservers(event)
			case observer, ok = <-distributor.toRegister:
				distributor.registerObserver(observer)
			case observer, ok = <-distributor.toUnregister:
				distributor.unregisterObserver(observer)
			case _, ok = <-stopObserving:
				ok = false
			}
		}
	}()

	return stopObserving
}

// RegisterObserver registers an observer to receive
// notifications with key found in keys, provided as
// second parameter.
func (distributor Distributor) RegisterObserver(observer Observer, keys ...string) {
	distributor.toRegister <- Observe{observer, keys}
}

// UnregisterObserver removes an observern from the
// list receiving notifications for keys in second
// parameter.
func (distributor Distributor) UnregisterObserver(observer Observer, keys ...string) {
	distributor.toUnregister <- Observe{observer, keys}
}

// RegisterObservers registers multiple observers to receive
// notifications with key found in Observe-object.
func (distributor Distributor) RegisterObservers(observe ...Observe) {
	for _, o := range observe {
		distributor.toRegister <- o
	}
}

// NotifyObservers distributes multiple events to observers
// listening for the same key.
func (distributor Distributor) NotifyObservers(events ...Event) {
	for _, event := range events {
		distributor.events <- event
	}
}

func (distributor Distributor) registerObserver(observer Observe) {
	for _, key := range observer.keys {
		distributor.observers[key] = append(distributor.observers[key], observer)
	}
}

func (distributor Distributor) unregisterObserver(observer Observe) {
	for _, key := range observer.keys {
		i := -1
		observers, ok := distributor.observers[key]
		if !ok {
			continue
		}

		for j, o := range observers {
			if o.UID() == observer.Observer.UID() {
				i = j
				break
			}
		}

		if i >= 0 {
			distributor.observers[key][i] = distributor.observers[key][len(distributor.observers[key])-1]
			distributor.observers[key] = distributor.observers[key][:len(distributor.observers[key])-1]
		}

		if len(distributor.observers[key]) == 0 {
			delete(distributor.observers, key)
		}
	}
}

func (distributor Distributor) notifyObservers(event Event) error {
	key, err := event.GetKey()
	if err != nil {
		return err
	}

	for _, observer := range distributor.observers[key] {
		if err := observer.HandleEvent(event); err != nil {
			return err
		}
	}

	return nil
}
