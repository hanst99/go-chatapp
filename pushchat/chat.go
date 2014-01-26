package pushchat

import (
    "net/http"
    "html/template"
    "encoding/json"
    "code.google.com/p/go.net/websocket"
    "sync"
    "io"
    "fmt"
)

type message struct {
    From    string
    Content string
}

type subscriber chan message

type chatRoom struct {
    Name        string
    messages    chan message
    subscribers []subscriber
    mutex       sync.Mutex
}

func NewRoom(name string) *chatRoom {
    return &chatRoom {
        Name: name,
        messages: make(chan message),
        subscribers: []subscriber{},
    }
}

func (this *chatRoom) AddSubscriber(sub subscriber) {
    this.mutex.Lock()
    defer this.mutex.Unlock()
    this.subscribers = append(this.subscribers,sub)
}

func (this *chatRoom) DispatchMessages() {
    for msg := range this.messages {
        this.mutex.Lock()
        for _,sub := range this.subscribers {
            sub<- msg
        }
        this.mutex.Unlock()
    }
}

func readMessages(r io.Reader, m chan message){
    decoder := json.NewDecoder(r)
    for {
        var msg message
        err := decoder.Decode(&msg)
        if err != nil {
            fmt.Printf("Error: %s",err)
            break
        }
        fmt.Printf("From: %s; Content: %s\n",msg.From,msg.Content)
        m<- msg
    }
}

func handleChatter(room *chatRoom, conn *websocket.Conn) {
    fmt.Println("Started handling a connection")
    encoder := json.NewEncoder(conn)
    incMessages := make(chan message)
    go readMessages(conn,room.messages)
    room.AddSubscriber(subscriber(incMessages))
    for msg := range incMessages {
        encoder.Encode(msg)
    }
}

func serveStaticFiles(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w,r,r.URL.Path[1:])
}



func StartApp(port uint16) error {
    defaultRoom := NewRoom("DefaultRoom")
    go defaultRoom.DispatchMessages()
    templ,err := template.ParseFiles("views/pushchat/index.html")
    if err != nil {
        return err
    }
    handleIndex := func (w http.ResponseWriter, r *http.Request) {
        templ.ExecuteTemplate(w,"push_chat",defaultRoom)
    }
    handler := http.NewServeMux()
    handler.HandleFunc("/chat",websocket.Handler(func(conn *websocket.Conn) {
        handleChatter(defaultRoom,conn)
    }).ServeHTTP)
    handler.HandleFunc("/public/",serveStaticFiles)
    handler.HandleFunc("/",handleIndex)
    http.ListenAndServe(fmt.Sprint(":",port),handler)
    return nil
}
