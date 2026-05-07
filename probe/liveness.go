package probe

import (
	"fmt"
	"net"
)

func Liveness(opts Options) error {
	conn, err := net.DialTimeout("tcp", opts.addr(), opts.Timeout)
	if err != nil {
		return fmt.Errorf("liveness failed: %w", err)
	}
	return conn.Close()
}
