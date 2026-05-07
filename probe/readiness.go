package probe

import (
	"errors"
	"fmt"

	"junhyung.kr/mc-probe/ping"
)

func Readiness(opts Options) error {
	resp, err := ping.Query(opts.addr(), opts.Timeout)
	if err != nil {
		return fmt.Errorf("readiness failed: %w", err)
	}
	if resp.Players.Max == 0 {
		return errors.New("server still bootstrapping")
	}
	if resp.Players.Online >= resp.Players.Max {
		return fmt.Errorf("server full (%d/%d)", resp.Players.Online, resp.Players.Max)
	}
	return nil
}
