//go:build !windows
// +build !windows

package services

import (
	"fmt"
)

type Window struct{}

func (w *Window) DisableConsoleQuickEdit() {
	fmt.Println("Not implemented for non-Windows systems")
}
