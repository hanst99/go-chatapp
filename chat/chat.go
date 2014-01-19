package chat;

import (
    "net/http"
    "fmt"
    "time"
    "github.com/hanst99/chatapp/web"
)

var sstorage *SessionStorage

func index(w http.ResponseWriter, r *http.Request) {
}

func startApp(port ushort) error {
    sstorage = CreateSessionStorage(SessionConfig { ValidFor: 15 * time.Minute })
    fmt.Printf("Created session storage\n")
}
