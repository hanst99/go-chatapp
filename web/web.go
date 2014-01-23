//Provides generic low level web functionality, like sessions
package web;

import (
    "net/http"
    "time"
    "strconv"
    "errors"
    "fmt"
    "sync"
)

//HTTP cookie based session
type Session struct {
    lastAccessed time.Time
    storage      map[string]string
}

func (this *Session) GetVal(name string)(string,error) {
    val,ok := this.storage[name]
    if !ok {
        return "", errors.New("No such value in session!")
    }
    return val,nil
}

func (this *Session) PutVal(name string, val string) {
    this.storage[name]=val
}

type SessionConfig struct {
    ValidFor time.Duration
    ClearInterval time.Duration
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
    storage.StartClearingSessions()
    return storage
}

// start clearing invalid (outdated) sessions from the storage,
// after every config.ClearInterval
//Note: Normally there's no need to call this function -
// it is done automatically upon creating a session storage
func (this *SessionStorage) StartClearingSessions() {
    go func() {
        for {
            select {
            case <-this.done:
                break
            default:
                this.cleanUp()
                time.Sleep(this.config.ClearInterval)
            }
        }
    }()
}

//Call this to stop clearing out of date sessions
//you should not use the session storage after calling this function
//Note: if StartClearingSessions was called more than once,
//this needs to be called once for each of those calls
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

//Creates a session and places its id in the session cookie of the response
//Note: Normally, you'd call GetSession instead
func (this *SessionStorage) CreateSession(w http.ResponseWriter) (*Session) {
        session := &Session{ lastAccessed: time.Now(), storage: make(map[string]string) }
        this.sessions[this.sessionCounter] = session
        http.SetCookie(w, &http.Cookie{Name: "session", Value: fmt.Sprint(this.sessionCounter)})
        this.sessionCounter += 1
        session.lastAccessed = time.Now()
        return session
}

//gets a session for this request
//if no session exists yet, create one
func (this *SessionStorage) GetSession(w http.ResponseWriter, req *http.Request) (*Session,error) {
    //look for session cookie
    sCookie,err := req.Cookie("session")
    var sessionId uint64
    this.mutex.Lock()
    defer this.mutex.Unlock()
    if(err != nil) {
        //if there's no session associated
        return this.CreateSession(w),nil
    } else {
        //if a session already exists
        sessionId,err = strconv.ParseUint(sCookie.Value,10,64)
        if(err != nil) {
            //if conversion failed
            return nil,err;
        }
    }
    session,ok := this.sessions[sessionId]
    if !ok {
        return this.CreateSession(w),nil
    }
    session.lastAccessed = time.Now()
    return session,nil
}
