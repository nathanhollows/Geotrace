package flash

import (
	"encoding/gob"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

func init() {
	gob.Register(Message{})
}

// Message is a struct containing each flashed message
type Message struct {
	Title   string
	Message string
	Style   string
}

func getCookieStore() *sessions.CookieStore {
	key := []byte(os.Getenv("GEOTRACE_SESSION_KEY"))
	return sessions.NewCookieStore(key)
}

// Set adds a new message into the cookie storage.
func Set(w http.ResponseWriter, r *http.Request, message Message) {
	store := getCookieStore()
	session, _ := store.Get(r, "om")
	session.Options.HttpOnly = true
	session.Options.Secure = true
	session.Options.SameSite = http.SameSiteStrictMode
	session.AddFlash(message)
	session.Save(r, w)
}

// Get flash messages from the cookie storage.
func Get(w http.ResponseWriter, r *http.Request) []interface{} {
	session, err := getCookieStore().Get(r, "om")
	if err == nil {
		messages := session.Flashes()
		if len(messages) > 0 {
			session.Save(r, w)
		}
		return messages
	}
	return nil
}
