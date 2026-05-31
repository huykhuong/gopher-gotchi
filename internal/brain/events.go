package brain

type EventType int

const (
	FileSaved EventType = iota
	CommandReceived
	SystemAlert
	ClipboardErrorDetected
)

type DataEvent struct {
	Type    EventType
	Payload any
}
