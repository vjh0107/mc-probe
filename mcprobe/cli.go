// Package mcprobe exposes mc-probe's CLI as a callable entrypoint so
// downstream binaries (e.g. furnace's cmd/mcprobe) can embed it without
// duplicating the kong wiring.
package mcprobe

import (
	"fmt"
	"time"

	"github.com/alecthomas/kong"

	"junhyung.kr/mc-probe/probe"
)

type CLI struct {
	Liveness  LivenessCmd      `cmd:"" help:"Liveness probe — server accepts TCP connections."`
	Readiness ReadinessCmd     `cmd:"" help:"Readiness probe — server has player slots available."`
	Startup   StartupCmd       `cmd:"" help:"Startup probe — server has finished bootstrapping."`
	Version   kong.VersionFlag `name:"version" help:"Print version information and exit."`
}

type probeFlags struct {
	Host    string        `help:"Server host." default:"127.0.0.1"`
	Port    int           `help:"Server port." default:"25565"`
	Timeout time.Duration `help:"Ping timeout." default:"1s"`
}

func (f probeFlags) options() probe.Options {
	return probe.Options{Host: f.Host, Port: f.Port, Timeout: f.Timeout}
}

type LivenessCmd struct{ probeFlags }

func (c *LivenessCmd) Run() error { return probe.Liveness(c.options()) }

type ReadinessCmd struct{ probeFlags }

func (c *ReadinessCmd) Run() error { return probe.Readiness(c.options()) }

type StartupCmd struct{ probeFlags }

func (c *StartupCmd) Run() error { return probe.Startup(c.options()) }

// Run parses os.Args and dispatches to the selected probe sub-command.
// version / buildTime are baked in by the caller via -ldflags '-X'.
func Run(version, buildTime string) {
	var cli CLI
	ctx := kong.Parse(&cli,
		kong.Name("mcprobe"),
		kong.Description("Minecraft server probe utility."),
		kong.Vars{"version": fmt.Sprintf("%s (built %s)", version, buildTime)},
	)
	ctx.FatalIfErrorf(ctx.Run())
}
