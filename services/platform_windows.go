// Author: Bruce Lu
// Email: lzbgt_AT_icloud.com

//go:build windows
// +build windows

package services

import (
	"fmt"

	"golang.org/x/sys/windows"
)

func (w *Window) DisableConsoleQuickEdit() {
	kernel32 := windows.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("SetConsoleMode")

	handle, err := windows.GetStdHandle(windows.STD_INPUT_HANDLE)
	if err != nil {
		fmt.Println("Error getting console handle:", err)
		return
	}

	var mode uint32
	err = windows.GetConsoleMode(handle, &mode)
	if err != nil {
		fmt.Println("Error getting console mode:", err)
		return
	}

	// Clear the ENABLE_QUICK_EDIT_MODE and ENABLE_INSERT_MODE flags
	mode &^= windows.ENABLE_QUICK_EDIT_MODE | windows.ENABLE_INSERT_MODE

	// Set the new mode
	r, _, err := proc.Call(uintptr(handle), uintptr(mode))
	if r == 0 {
		fmt.Println("Error setting console mode:", err)
		return
	}

}
