package logic

// Notification represents a single notification/toast
type Notification struct {
	ID         string           `json:"id"`
	Type       NotificationType `json:"type"`
	Title      string           `json:"title"`
	Message    string           `json:"message"`
	Duration   int              `json:"duration"`
	Cancelable bool             `json:"cancelable"`
}

// NotificationType defines the kind of notification
type NotificationType string

const (
	NotificationTypeSuccess NotificationType = "success"
	NotificationTypeInfo    NotificationType = "info"
	NotificationTypeWarning NotificationType = "warning"
	NotificationTypeError   NotificationType = "error"
)
