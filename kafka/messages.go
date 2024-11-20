package kafka

import (
	"encoding/json"
)

type NotificationMessage struct {
	Username       string          `json:"username"`
	NotificationID string          `json:"notificationID"`
	Priority       MessagePriority `json:"priority"`
}

func (e *NotificationMessage) GetPriority() MessagePriority {
	return e.Priority
}

func (e *NotificationMessage) Encode() ([]byte, error) {
	return json.Marshal(e)
}
