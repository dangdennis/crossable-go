package messages

type MessageType string

// String returns the string type of MessageType
func (m MessageType) String() string {
	return string(m)
}

const (
	MessageTypeIntro      MessageType = "intro"
	MessageTypeCompletion MessageType = "completion"
)
