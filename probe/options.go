package probe

import (
	"net"
	"strconv"
	"time"
)

type Options struct {
	Host    string
	Port    int
	Timeout time.Duration
}

func (o Options) addr() string {
	return net.JoinHostPort(o.Host, strconv.Itoa(o.Port))
}
