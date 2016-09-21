package observer

// Event interface provides methods to
// identify observers to distribute event
// to, based on a key.
// Instantiated objects following the Event will
// be distributed to observers event handlers.
type Event interface {
	GetKey() (string, error)
}
