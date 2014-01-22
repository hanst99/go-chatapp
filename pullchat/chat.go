package pullchat;

import (
    "net/http"
    "fmt"
    "time"
    "github.com/hanst99/chatapp/web"
    "html/template"
    "log"
)

var sstorage *web.SessionStorage

type user struct {
    Name string
    color uint32
}

type message struct {
    from user
    content string
}

type chatRoom struct {
    messages []message
    name string
}

func index(indexTemplate *template.Template, w http.ResponseWriter, r *http.Request) {
    err := indexTemplate.Execute(w,user {Name: "hannes"})
    if err != nil {
       log.Fatal(err)
    }
}

type TemplateServer func(*template.Template, http.ResponseWriter, *http.Request)

func wrapTemplate(templ *template.Template, server TemplateServer) http.HandlerFunc {
    return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
        server(templ,w,r)
    })
}

func StartApp(port uint16) error {
    sstorage = web.CreateSessionStorage(web.SessionConfig { ValidFor: 15 * time.Minute })
    indexTemplate,err := template.ParseFiles("views/pullchat/index.html")
    if err != nil {
        return err
    }
    http.HandleFunc("/",wrapTemplate(indexTemplate,index))
    err = http.ListenAndServe(fmt.Sprint(":",port), nil)
    if err != nil {
        return err
    }
    return nil
}
