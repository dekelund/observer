package observer

// Observer interface specifies
// handlers for notification-events
// distributed by observer package.
type Observer interface {
	HandleEvent(Event) error
}

// Observe provides mapping between
// observer and keys, usable when
// registrering observers.
type Observe struct {
	Observer
	keys []string
}
