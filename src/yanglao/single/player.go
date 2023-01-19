package single

import (
	"net/http"
	"sync"
	"yanglao/gonet"

	"yanglao/hcc/structure"

	goutils "yanglao/gonet/utils"

	"github.com/cihub/seelog"
)

const (
	PowerAdmin = "admin"
	PowerGuest = "guest"
)

type Player struct {
	Account  string `bson:"_id"`
	Passwd   string `bson:"pswd"`
	UserName string
	Phone    string
	Session  string `bson:"-"`
	Power    uint32
}

const (
	TypeAccount = iota
	TypeSession
)

type Players struct {
	lock        sync.Mutex
	players_acc map[string]*Player
}

func (p *Players) Init() {
	p.players_acc = make(map[string]*Player)
}

func (p *Players) Add(player *Player) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.players_acc[player.Account] = player
}

func (p *Players) Get(account string) *Player {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.players_acc[account]
}

type PlayerManager struct {
	players Players
}

func (p *PlayerManager) Init() {
	p.players.Init()
}

func (p *PlayerManager) Get(account string) *Player {
	return p.players.Get(account)
}

func (p *PlayerManager) GetBySessionstr(sessionstr string) *Player {
	session := SessionMgr.GetSession(sessionstr)
	if session == nil {
		return nil
	}

	account := session.Get(SessionKey_Acc)
	if account == nil || account == "" {
		return nil
	}
	return p.Get(account.(string))
}

func (p *PlayerManager) GetByRequest(r *http.Request) *Player {
	session := SessionMgr.GetByRequest(r)
	seelog.Info("GetByRequest session", session)
	if session == nil {
		return nil
	}

	account := session.Get(SessionKey_Acc)
	if account == nil || account == "" {
		return nil
	}
	return p.Get(account.(string))
}

func (p *PlayerManager) Load(account, pwd string) *Player {
	player := p.Get(account)
	if player != nil && player.Account == account {
		if player.Passwd == pwd {
			return player
		}
		return nil
	}

	ret, err := gonet.CallByName("HccMysqlSvr", "Login", account, pwd)
	if err != nil {
		panic(err)
	}

	user := &structure.User{}
	err = goutils.ExpandResult(ret, &user)
	if err != nil {
		panic(err)
	}
	if user.Account == "" {
		return nil
	}

	dbPlayer := Player{
		Account:  user.Account,
		Passwd:   user.PassWord,
		UserName: user.UserName,
		Phone:    user.Phone,
		Power:    user.Power,
		Session:  ""}
	p.players.Add(&dbPlayer)
	return &dbPlayer
}

var PlayerMgr *PlayerManager

func NewPlayerManager() *PlayerManager {
	mgr := new(PlayerManager)
	mgr.Init()
	return mgr
}

//func init() {
//	PlayerMgr = NewPlayerManager()
//}
