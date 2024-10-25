//go:build !windows
// +build !windows

package services

import (
	"fmt"
)

func (w *Window) DisableConsoleQuickEdit() {
	fmt.Println("Not implemented for non-Windows systems")
}
