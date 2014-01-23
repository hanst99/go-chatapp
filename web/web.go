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

func (this *Session) GetVal(name string)(string,error) {

}

type SessionConfig struct {
    ValidFor time.Duration
}

type SessionStorage struct {
    sessions       map[uint64]*Session
    sessionCounter uint64
    config         SessionConfig
    done           chan bool
    mutex          sync.Mutex
}


//Create a new session storage, using the given config
func CreateSessionStorage(config SessionConfig) *SessionStorage {
    storage := &SessionStorage {
        sessions: make(map[uint64]*Session),
        sessionCounter: 0,
        done: make(chan bool),
        config: config,
    }
    //start clearing sessions every 5 minutes
    //until told to stop
    go func() {
        for {
            select {
            case <-storage.done:
                break
            default:
                storage.cleanUp()
                time.Sleep(5 * time.Minute)
            }
        }
    }()
    return storage
}

//Call this to stop clearing out of date sessions
//you should not use the session storage after calling this function
func (this *SessionStorage) StopClearingSessions() {
    this.done<-true
}

//clears expires sessions
func (this *SessionStorage) cleanUp() {
    this.mutex.Lock()
    defer this.mutex.Unlock()
    for sessionId,session := range this.sessions {
        if(time.Since(session.lastAccessed) > this.config.ValidFor) {
            delete(this.sessions,sessionId)
        }
    }
}


//gets a session for this request
//if no session exists yet, create one
func (this *SessionStorage) GetSession(req *http.Request) (*Session,error) {
    //look for session cookie
    sCookie,err := req.Cookie("session")
    var sessionId uint64
    this.mutex.Lock()
    defer this.mutex.Unlock()
    if(err != nil) {
        //if there's no session associated
        this.sessions[this.sessionCounter] = &Session{ lastAccessed: time.Now(), storage: make(map[string]string) };
        sessionId = this.sessionCounter
        req.AddCookie(&http.Cookie{Name: "session", Value: fmt.Sprint(sessionId)})
        this.sessionCounter += 1
    } else {
        //if a session already exists
        sessionId,err = strconv.ParseUint(sCookie.Value,10,64)
        if(err != nil) {
            //if conversion failed
            return nil,err;
        }
    }
    session := this.sessions[sessionId]
    session.lastAccessed = time.Now()
    return session,nil
}
