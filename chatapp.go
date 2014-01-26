package main;

import (
    "github.com/hanst99/go-chatapp/pullchat"
    "github.com/hanst99/go-chatapp/pushchat"
    "flag"
    "fmt"
)

func main() {
    var port int
    flag.IntVar(&port,"port",0,"The port to start the web app on")
    var mode string
    flag.StringVar(&mode,"mode","pull","Mode of the application (push|pull)")
    flag.Parse()
    if port <= 0  {
        panic(fmt.Errorf("Must give valid port (passed %d)",port))
    }
    switch mode {
    case "pull":
        err := pullchat.StartApp(uint16(port))
        if err != nil {
            panic(err)
        }
    case "push":
        err := pushchat.StartApp(uint16(port))
        if err != nil {
            panic(err)
        }
    default:
        panic(fmt.Errorf("No such mode: %s",mode))
    }
}
