package logic

import (
	"fmt"
	"sync"
)

// notifList is the in-memory notification store
var (
	notifList = make([]Notification, 0)
	notifID   = 0
	mu        sync.RWMutex
)

// AddNotification adds a notification to the store (no channels, no goroutines)
func AddNotification(notif Notification) {
	mu.Lock()
	defer mu.Unlock()
	notif.ID = fmt.Sprintf("notif_%d", notifID+1)
	notifID++
	notifList = append(notifList, notif)
}

// GetNotifications returns a copy of the notification list
func GetNotifications() []Notification {
	mu.RLock()
	defer mu.RUnlock()
	resp := make([]Notification, len(notifList))
	copy(resp, notifList)
	return resp
}

// RemoveNotification removes a notification by ID
func RemoveNotification(id string) {
	mu.Lock()
	defer mu.Unlock()
	for i, n := range notifList {
		if n.ID == id {
			notifList = append(notifList[:i], notifList[i+1:]...)
			break
		}
	}
}

// ClearNotifications clears all notifications
func ClearNotifications() {
	mu.Lock()
	defer mu.Unlock()
	notifList = make([]Notification, 0)
}
