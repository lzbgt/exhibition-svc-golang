// Author: Bruce Lu
// Email: lzbgt_AT_icloud.com

//go:build !windows
// +build !windows

package services

import (
	"fmt"
)

func (w *Window) DisableConsoleQuickEdit() {
	fmt.Println("Not implemented for non-Windows systems")
}
