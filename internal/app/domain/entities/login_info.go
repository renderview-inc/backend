package entities

import (
	"net"
	"time"

	"github.com/google/uuid"
)

type LoginInfo struct {
	id        uuid.UUID
	userID    uuid.UUID
	loginTime time.Time
	userAgent string
	ipAddr    net.IP
	success   bool
}

func NewLoginInfo(id uuid.UUID, userID uuid.UUID, loginTime time.Time,
	userAgent string, ipAddr net.IP, success bool) LoginInfo {
	return LoginInfo{
		id, userID, loginTime, userAgent, ipAddr, success,
	}
}

func (li *LoginInfo) ID() uuid.UUID {
	return li.id
}

func (li *LoginInfo) UserID() uuid.UUID {
	return li.userID
}

func (li *LoginInfo) LoginTime() time.Time {
	return li.loginTime
}

func (li *LoginInfo) UserAgent() string {
	return li.userAgent
}

func (li *LoginInfo) IpAddr() net.IP {
	return li.ipAddr
}

func (li *LoginInfo) Success() bool {
	return li.success
}
