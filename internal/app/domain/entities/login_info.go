package entities

import (
	"net"
	"time"

	"github.com/google/uuid"
)

type LoginInfo struct {
	Id        uuid.UUID
	UserID   uuid.UUID
	LoginTime time.Time
	UserAgent string
	IpAddr    net.IP
	Success   bool
}

func NewLoginInfo(id uuid.UUID, userID uuid.UUID, loginTime time.Time,
	userAgent string, ipAddr net.IP, success bool) LoginInfo {
	return LoginInfo{
		id, userID, loginTime, userAgent, ipAddr, success,
	}
}

func (li *LoginInfo) GetID() uuid.UUID {
	return li.Id
}

func (li *LoginInfo) UserGetID() uuid.UUID {
	return li.UserID
}

func (li *LoginInfo) GetLoginTime() time.Time {
	return li.LoginTime
}

func (li *LoginInfo) GetUserAgent() string {
	return li.UserAgent
}

func (li *LoginInfo) GetIpAddr() net.IP {
	return li.IpAddr
}

func (li *LoginInfo) GetSuccess() bool {
	return li.Success
}
