package single

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"yanglao/static"

	"github.com/cihub/seelog"
)

var (
	SessionKey_Acc = "account"
)

type Session interface {
	Set(key, value interface{}) error //设置Session
	Get(key interface{}) interface{}  //获取Session
	Delete(key interface{}) error     //删除Session
	SessionID() string                //当前SessionID
	SetLast(l time.Time)
	Last() time.Time
}

type Session4Memory struct {
	sid  string
	lock sync.Mutex
	last time.Time
	data map[interface{}]interface{}
}

func newSession4Memory() *Session4Memory {
	return &Session4Memory{
		data: make(map[interface{}]interface{}),
	}
}

var (
	err_session_exit error = errors.New("session exit")
)

//同一个会话均可调用，进行设置，改操作必须拥有排斥锁
func (si *Session4Memory) Set(key, value interface{}) error {
	si.lock.Lock()
	defer si.lock.Unlock()

	if _, ok := si.data[key]; ok {
		return err_session_exit
	}

	si.data[key] = value
	return nil
}

func (si *Session4Memory) Get(key interface{}) interface{} {
	if value := si.data[key]; value != nil {
		return value
	}
	return nil
}
func (si *Session4Memory) Delete(key interface{}) error {
	if value := si.data[key]; value != nil {
		delete(si.data, key)
	}
	return nil
}
func (si *Session4Memory) SessionID() string {
	return si.sid
}

func (si *Session4Memory) SetLast(l time.Time) {
	si.last = l
}

func (si *Session4Memory) Last() time.Time {
	return si.last
}

type Storage interface {
	InitSession(sid string) (Session, error)
	SetSession(session Session) error
	GetSession(sid string) Session
	DestroySession(sid string) error
	GCSession()

	GetMaxLife() int64
	ShowSessions()
}

type SessionMemoryMgr struct {
	lock     sync.Mutex //互斥锁
	life     int64      //超时时间
	sessions map[string]Session
}

func newFromMemory() *SessionMemoryMgr {
	return &SessionMemoryMgr{
		life: static.HttpConfig.CookieLife, //60 * 30 * 1000,
		//life:     60 * 1 * 1000,
		sessions: make(map[string]Session, 0),
	}
}

//初始换会话session，这个结构体操作实现Session接口
func (fm *SessionMemoryMgr) InitSession(sid string) (Session, error) {
	fm.lock.Lock()
	defer fm.lock.Unlock()

	newSession := newSession4Memory()
	newSession.sid = sid
	newSession.last = time.Now()

	fm.sessions[sid] = newSession //内存管理map
	//seelog.Info("InitSession:", sid, fm.sessions)
	return newSession, nil
}

//设置
func (fm *SessionMemoryMgr) SetSession(session Session) error {
	fm.sessions[session.SessionID()] = session
	return nil
}

//获取
func (fm *SessionMemoryMgr) GetSession(sid string) Session {
	fmt.Println(fm.sessions)
	session, ok := fm.sessions[sid]
	if !ok {
		return nil
	}
	session.SetLast(time.Now())
	return session
}

//销毁session
func (fm *SessionMemoryMgr) DestroySession(sid string) error {
	if _, ok := fm.sessions[sid]; ok {
		delete(fm.sessions, sid)
		return nil
	}
	return nil
}

//监判超时
func (fm *SessionMemoryMgr) GCSession() {

	sessions := fm.sessions

	if len(sessions) < 1 {
		return
	}

	for k, v := range sessions {
		t := (v.(*Session4Memory).last.Unix()) + (fm.life / 1000)

		if t <= time.Now().Unix() { //超时了
			delete(fm.sessions, k)
			seelog.Info("GCSession", k)
		}
	}

}

func (fm *SessionMemoryMgr) GetMaxLife() int64 {
	return fm.life
}

func (fm *SessionMemoryMgr) ShowSessions() {
	fmt.Println(fm.sessions)
}

///////////////////////////////////////////////////////////////////

type SessionManager struct {
	cookieName string
	storage    Storage
	lock       sync.Mutex
	lockGC     sync.Mutex
	lockR      sync.Mutex
}

func NewSessionManager() *SessionManager {
	sessionMgr := &SessionManager{
		cookieName: static.HttpConfig.CookieName, //"bm-cookie",
		storage:    newFromMemory(),
	}
	go sessionMgr.GC()
	return sessionMgr
}

//先判断当前请求的cookie中是否存在有效的session,存在返回，不存在创建
func (m *SessionManager) BeginSession(w http.ResponseWriter, r *http.Request) Session {
	//防止处理时，进入另外的请求
	m.lock.Lock()
	defer m.lock.Unlock()

	cookie, err := r.Cookie(m.cookieName)
	if err != nil || cookie.Value == "" { //如果当前请求没有改cookie名字对应的cookie
		//创建一个
		sid := m.randomId()
		//根据保存session方式，如内存，数据库中创建
		session, _ := m.storage.InitSession(sid) //该方法有自己的锁，多处调用到

		//用session的ID于cookie关联
		//cookie名字和失效时间由session管理器维护
		cookie := http.Cookie{
			Name: m.cookieName,
			//这里是并发不安全的，但是这个方法已上锁
			Value:    url.QueryEscape(sid), //转义特殊符号@#￥%+*-等
			Path:     "/",
			Domain:   static.HttpConfig.Domain, //"localhost",
			HttpOnly: false,
			MaxAge:   int(m.storage.GetMaxLife() / 1000),
			Expires:  session.Last().Add(time.Millisecond * time.Duration(m.storage.GetMaxLife())),
		}
		http.SetCookie(w, &cookie) //设置到响应中
		return session
	} else { //如果存在
		sid, _ := url.QueryUnescape(cookie.Value)              //反转义特殊符号
		session := m.storage.(*SessionMemoryMgr).sessions[sid] //从保存session介质中获取
		if session == nil {
			//创建一个
			sid := m.randomId()
			//根据保存session方式，如内存，数据库中创建
			newSession, _ := m.storage.InitSession(sid) //该方法有自己的锁，多处调用到

			//用session的ID于cookie关联
			//cookie名字和失效时间由session管理器维护
			newCookie := http.Cookie{
				Name: m.cookieName,
				//这里是并发不安全的，但是这个方法已上锁
				Value:    url.QueryEscape(sid), //转义特殊符号@#￥%+*-等
				Path:     "/",
				Domain:   static.HttpConfig.Domain, //"localhost",
				HttpOnly: false,
				MaxAge:   int(m.storage.GetMaxLife() / 1000),
				Expires:  newSession.Last().Add(time.Millisecond * time.Duration(m.storage.GetMaxLife())),
			}
			http.SetCookie(w, &newCookie) //设置到响应中
			return newSession
		}
		return session
	}

}

func (m *SessionManager) EndSession(session string) {
	fmt.Println("end session", session)
	m.lock.Lock()
	defer m.lock.Unlock()
	m.storage.DestroySession(session)

	m.storage.ShowSessions()
}

//开启每个会话，同时定时调用该方法
//到达session最大生命时，且超时时。回收它
func (m *SessionManager) GC() {
	m.lockGC.Lock()
	defer m.lockGC.Unlock()

	m.storage.GCSession()
	//在多长时间后执行匿名函数，这里指在某个时间后执行GC
	time.AfterFunc(time.Duration(m.storage.GetMaxLife()*1000*1000/10), func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05.000"), "SessionManager GC ", m.storage.GetMaxLife())
		m.GC()
	})
}

//生成一定长度的随机数
func (m *SessionManager) randomId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	//加密
	return base64.URLEncoding.EncodeToString(b)
}

func (m *SessionManager) GetByRequest(r *http.Request) Session {
	cookie, err := r.Cookie(m.cookieName)
	if err != nil || cookie.Value == "" {
		return nil
	}
	se, err := url.QueryUnescape(cookie.Value)
	return m.storage.GetSession(se)
}

func (m *SessionManager) GetSession(str string) Session {
	//m.storage.ShowSessions()
	return m.storage.GetSession(str)
}

func (m *SessionManager) SetCookie(w http.ResponseWriter, sessionid string) {
	m.lockR.Lock()
	defer m.lockR.Unlock()

	session := m.GetSession(sessionid)
	if session == nil {
		seelog.Error("SessionManager::SetCookie not found session by id", sessionid)
		return
	}
	cookie := http.Cookie{
		Name:     m.cookieName,
		Value:    url.QueryEscape(session.SessionID()), //转义特殊符号@#￥%+*-等
		Path:     "/",
		Domain:   static.HttpConfig.Domain, //"localhost",
		HttpOnly: false,
		MaxAge:   int(m.storage.GetMaxLife() / 1000),
		Expires:  session.Last().Add(time.Millisecond * time.Duration(m.storage.GetMaxLife())),
	}
	http.SetCookie(w, &cookie) //设置到响应中
}

var SessionMgr *SessionManager

//func init() {
//	SessionMgr = NewSessionManager()
//}
