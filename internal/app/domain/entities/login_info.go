package entities

import (
	"net"
	"time"

	"github.com/google/uuid"
)

type LoginInfo struct {
	Id        uuid.UUID
	UserID    uuid.UUID
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
