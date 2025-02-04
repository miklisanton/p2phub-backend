package utils

import (
	"github.com/rs/zerolog"
	"runtime"
	"strconv"
	"strings"
)

func GetGoroutineID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	stack := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, _ := strconv.Atoi(stack)
	return id
}

// GoroutineHook adds the goroutine ID to log events
type GoroutineHook struct{}

// Run implements the zerolog.Hook interface
func (h GoroutineHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	e.Int("gid", GetGoroutineID())
}
