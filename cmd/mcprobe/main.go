package main

import "junhyung.kr/mc-probe/mcprobe"

var (
	version   = "dev"
	buildTime = "unknown"
)

func main() { mcprobe.Run(version, buildTime) }
