package observer

// Observer interface specifies
// handlers for notification-events
// distributed by observer package.
// UID function must return an uniq
// id for this client.
type Observer interface {
	HandleEvent(Event) error
	UID() string
}

// Observe provides mapping between
// observer and keys, usable when
// registrering observers.
type Observe struct {
	Observer
	keys []string
}
