package dtos

import "net"

type FullLoginInfo struct {
	Credentials Credentials `json:"credentials"`
	LoginMeta   LoginMeta   `json:"login_metadata"`
}

type Login struct {
	Credentials Credentials `json:"credentials"`
}

type LoginMeta struct {
	UserAgent string
	IpAddr    net.IP
}
