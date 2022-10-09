package main

import (
	"os"
	"syscall"
)

func redirect() {
	logFile, _ := os.OpenFile("panic.txt", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0755)
	syscall.Dup2(int(logFile.Fd()), 2)
}
