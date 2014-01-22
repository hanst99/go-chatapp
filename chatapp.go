package main;

import (
    "github.com/hanst99/go-chatapp/pullchat"
)

func main() {
    err := pullchat.StartApp(1234)
    if err != nil {
        panic(err)
    }
}
