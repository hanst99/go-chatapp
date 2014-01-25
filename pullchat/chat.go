package pullchat;

import (
    "net/http"
    "fmt"
    "time"
    "html/template"
    "github.com/hanst99/go-chatapp/web"
    "log"
    "sync"
    "io/ioutil"
)


type user struct {
    Name string
}

type message struct {
    From user
    At time.Time
    Content string
}

type chatRoom struct {
    Messages []message
    Name string
    mutex sync.Mutex
}

func (this *chatRoom) PostMessage(m message) {
    this.mutex.Lock()
    defer this.mutex.Unlock()
    this.Messages = append(this.Messages,m)
}


var sessionStorage *web.SessionStorage
var defaultRoom chatRoom = chatRoom{Messages: []message{},Name: "Default Room"}

func postMessage(w http.ResponseWriter, r *http.Request) {
    if(r.Method != "POST") {
        log.Fatal("/post_message must be POST'ed to!")
    }
    session,err := sessionStorage.GetSession(w,r)
    if err != nil {
        log.Fatal(err)
    }
    username,err := session.GetVal("user.name")
    if err != nil {
        log.Fatal(err)
    }
    content,err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Fatal(err)
    }
    defaultRoom.PostMessage(message {From: user{username}, At: time.Now(), Content: string(content)})
}

// the index page of the chat
func index(indexTemplate *template.Template, w http.ResponseWriter, r *http.Request) {
    //disable caching: HTTP/1.1 HTTP/1.0 and proxies
    //from http://stackoverflow.com/questions/49547/making-sure-a-web-page-is-not-cached-across-all-browsers
    w.Header().Set("Cache-Control","no-cache,no-store,must-revalidate")
    w.Header().Set("Pragma:", "no-cache")
    w.Header().Set("Expires","0")
    session,err := sessionStorage.GetSession(w,r)
    if err != nil {
        log.Fatal(err)
    }
    if r.Method == "POST" {
        user := r.FormValue("username")
        session.PutVal("user.name",user)
    }
    _,err = session.GetVal("user.name")
    if err != nil {
        err = indexTemplate.ExecuteTemplate(w,"signup",nil)
        if err != nil {
            log.Fatal(err)
        }
        return
    }
    defaultRoom.mutex.Lock()
    defer defaultRoom.mutex.Unlock()
    err = indexTemplate.ExecuteTemplate(w,"pull_chat",defaultRoom)
    if err != nil {
       log.Fatal(err)
    }
}

func public(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w,r,r.URL.Path[1:])
}

type templateServer func(*template.Template, http.ResponseWriter, *http.Request)

func wrapTemplate(templ *template.Template, server templateServer) http.HandlerFunc {
    return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
        server(templ,w,r)
    })
}

func StartApp(port uint16) error {
    sessionStorage = web.CreateSessionStorage(web.SessionConfig{
        ValidFor: 15*time.Minute,
        ClearInterval: 5*time.Minute})
    defer sessionStorage.StopClearingSessions()
    indexTemplate,err := template.ParseFiles("views/pullchat/index.html","views/pullchat/signup.html")
    if err != nil {
        return err
    }
    handler := http.NewServeMux()
    handler.HandleFunc("/",wrapTemplate(indexTemplate,index))
    handler.HandleFunc("/public/",public)
    handler.HandleFunc("/post_message",postMessage)
    err = http.ListenAndServe(fmt.Sprint(":",port), handler)
    if err != nil {
        return err
    }
    return nil
}
