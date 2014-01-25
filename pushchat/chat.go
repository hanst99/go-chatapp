package pushchat

import (
    "github.com/hanst99/go-chatapp/web"
    "fmt"
    "net/http"
    "time"
    "html/template"
)

type message struct {
    From    string
    At      time.Time
    Content string
}

type subscriber chan message

type chatRoom struct {
    name        string
    messages    chan message
    subscribers []subscriber
}

func (this *chatRoom) AddSubscriber(room *chatRoom, sub subscriber) {
    append(this.subscribers,sub)
}

var sessionStorage *web.SessionStorage


func StartApp(port uint16) error {
    sessionStorage = web.CreateSessionStorage(web.SessionConfig{
        ValidFor: 15*time.Minute,
        ClearInterval: 5*time.Minute})
    handler := http.NewServeMux()
    defer sessionStorage.StopClearingSessions()
    return nil
}
