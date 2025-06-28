package main

import (
	"github.com/core-tools/hsu-echo/cmd/srv/domain"
	"github.com/core-tools/hsu-echo/pkg/control"
)

func main() {
	control.MainEcho(domain.NewSimpleHandler)
}
