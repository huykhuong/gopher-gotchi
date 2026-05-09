package brain

type EventType int

const (
	FileSaved EventType = iota
	CommandReceived
	SystemAlert
)

type DataEvent struct {
	Type    EventType
	Payload any
}
