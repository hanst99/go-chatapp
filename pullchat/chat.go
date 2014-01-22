package pullchat;

import (
    "net/http"
    "fmt"
    "time"
    "html/template"
    "log"
)


type user struct {
    Name string
    Color uint32
}

type message struct {
    From user
    At time.Time
    Content string
}

type chatRoom struct {
    Messages []message
    Name string
}

func index(indexTemplate *template.Template, w http.ResponseWriter, r *http.Request) {
    //disable caching: HTTP/1.1 HTTP/1.0 and proxies
    //from http://stackoverflow.com/questions/49547/making-sure-a-web-page-is-not-cached-across-all-browsers
    w.Header().Set("Cache-Control","no-cache,no-store,must-revalidate")
    w.Header().Set("Pragma:", "no-cache")
    w.Header().Set("Expires","0")
    err := indexTemplate.Execute(w,user {Name: "hannes"})
    if err != nil {
       log.Fatal(err)
    }
}

func public(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w,r,r.URL.Path[1:])
}

type TemplateServer func(*template.Template, http.ResponseWriter, *http.Request)

func wrapTemplate(templ *template.Template, server TemplateServer) http.HandlerFunc {
    return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
        server(templ,w,r)
    })
}

func StartApp(port uint16) error {
    indexTemplate,err := template.ParseFiles("views/pullchat/index.html")
    if err != nil {
        return err
    }
    handler := http.NewServeMux()
    handler.HandleFunc("/",wrapTemplate(indexTemplate,index))
    handler.HandleFunc("/public/",public)
    err = http.ListenAndServe(fmt.Sprint(":",port), handler)
    if err != nil {
        return err
    }
    return nil
}
