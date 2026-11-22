// Package logic contiene funciones de parseo, validación y autenticación
// comunes usadas en los handlers.
package logic

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"
	"sync"
)

// NotificationType identifica la variante visual de una notificación.
type NotificationType string

const (
	NotificationSuccess NotificationType = "success"
	NotificationError   NotificationType = "error"
	NotificationInfo    NotificationType = "info"
	NotificationWarning NotificationType = "warning"
)

// Notification es un aviso efímero para el usuario.
type Notification struct {
	Type       NotificationType
	Label      string // etiqueta pequeña del encabezado (p.ej. "éxito")
	Title      string // título principal
	Message    string // detalle opcional
	ActionHref string // enlace opcional a la derecha
	ActionText string
}

// notificationEvent es un mensaje del canal: encola un aviso nuevo o pide
// (y consume) los pendientes de una sesión.
type notificationEvent struct {
	session      string
	notification *Notification
	resp         chan []Notification
}

var (
	notifyEvents = make(chan notificationEvent, 256)
	startOnce    sync.Once
)

// StartNotificationDispatcher arranca la gorutina única que atiende el canal.
func StartNotificationDispatcher() {
	startOnce.Do(func() {
		go func() {
			store := make(map[string][]Notification)
			for ev := range notifyEvents {
				if ev.resp != nil {
					items := store[ev.session]
					delete(store, ev.session)
					ev.resp <- items
					continue
				}
				store[ev.session] = append(store[ev.session], *ev.notification)
			}
		}()
	})
}

// AddNotification encola un aviso sin bloquear al productor.
func AddNotification(session string, n Notification) {
	if session == "" {
		session = "anon"
	}
	StartNotificationDispatcher()
	select {
	case notifyEvents <- notificationEvent{session: session, notification: &n}:
	default: // canal saturado: se descarta antes de bloquear
	}
}

// TakeNotifications devuelve y consume las notificaciones pendientes de la sesión.
func TakeNotifications(session string) []Notification {
	if session == "" {
		session = "anon"
	}
	StartNotificationDispatcher()
	resp := make(chan []Notification, 1)
	notifyEvents <- notificationEvent{session: session, resp: resp}
	return <-resp
}

// SessionID devuelve un identificador estable por navegador: usa user_id si hay
// sesión iniciada; si no, crea/reutiliza una cookie "sid" anónima.
func SessionID(w http.ResponseWriter, r *http.Request) string {
	if c, err := r.Cookie("user_id"); err == nil && c.Value != "" {
		return "user:" + c.Value
	}
	if c, err := r.Cookie("sid"); err == nil && c.Value != "" {
		return "sid:" + c.Value
	}
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "sid:" + strconv.FormatInt(r.ContentLength, 10)
	}
	sid := hex.EncodeToString(buf)
	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    sid,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
	return "sid:" + sid
}
