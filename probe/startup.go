package probe

import (
	"errors"
	"fmt"

	"junhyung.kr/mc-probe/ping"
)

func Startup(opts Options) error {
	resp, err := ping.Query(opts.addr(), opts.Timeout)
	if err != nil {
		return fmt.Errorf("startup failed: %w", err)
	}
	if resp.Players.Max == 0 {
		return errors.New("server still bootstrapping")
	}
	return nil
}
