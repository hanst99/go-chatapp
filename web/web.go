//Provides generic low level web functionality, like sessions
package web;

import (
    "net/http"
    "time"
    "strconv"
    "fmt"
    "sync"
)

//HTTP cookie based session
type Session struct {
    lastAccessed time.Time
    storage      map[string]string
}

type SessionConfig struct {
    ValidFor time.Duration
}

type SessionStorage struct {
    sessions       map[uint64]Session
    sessionCounter uint64
    config         SessionConfig
    mutex          sync.Mutex
}


func CreateSessionStorage(config SessionConfig) *SessionStorage {
    storage := &SessionStorage {
        sessions: make(map[uint64]Session),
        sessionCounter: 0,
        config: config,
    }
    go func() {
        for {
            storage.mutex.Lock()
            storage.cleanUp()
            storage.mutex.Unlock()
            time.Sleep(5 * time.Minute)
        }
    }()
    return storage
}

func (this *SessionStorage) cleanUp() {
    this.mutex.Lock()
    defer this.mutex.Unlock()
    now := time.Now()
    for sessionId,session := range this.sessions {
        if(now.Sub(session.lastAccessed) > this.config.ValidFor) {
            delete(this.sessions,sessionId)
        }
    }
}


//gets a session for this request
//if no session exists yet, create one
func (this *SessionStorage) GetSession(req *http.Request) (Session,error) {
    //look for session cookie
    sCookie,err := req.Cookie("session")
    var sessionId uint64
    this.mutex.Lock()
    defer this.mutex.Unlock()
    if(err != nil) {
        //if there's no session associated
        this.sessions[this.sessionCounter] = Session{ lastAccessed: time.Now(), storage: make(map[string]string) };
        sessionId = this.sessionCounter
        req.AddCookie(&http.Cookie{Name: "session", Value: fmt.Sprint(sessionId)})
        this.sessionCounter += 1
    } else {
        //if a session already exists
        sessionId,err = strconv.ParseUint(sCookie.Value,10,64)
        if(err != nil) {
            //if conversion failed
            return Session{},err;
        }
    }
    session := this.sessions[sessionId]
    session.lastAccessed = time.Now()
    return session,nil
}
