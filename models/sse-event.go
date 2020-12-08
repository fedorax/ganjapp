package models

// For streaming push data

// StreamChannel contains a map of user to channel
var StreamChannel = make(map[uint]chan SSEEvent)

// StreamStatus contains a map of user to channel status
var StreamStatus = make(map[uint]chan bool)

// StreamMessageID contains a map of user to message count
var StreamMessageID = make(map[uint]uint)

// SSEEvent is a model for server side events
type SSEEvent struct {
	Type string
	Data string
}

// CreateUserStream ensures that the passed user ID is present in the stream mappings
func CreateUserStream(userID uint) {
	if _, ok := StreamChannel[userID]; !ok {
		StreamChannel[userID] = make(chan SSEEvent)
		StreamStatus[userID] = make(chan bool)
		StreamMessageID[userID] = 0
	}
}
