package dtos

import "net"

type LoginMetaDto struct {
	UserAgent string
	IpAddr    net.IP
}
