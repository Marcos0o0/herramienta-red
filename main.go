package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/marcos0o0/herramienta-red/internal/diagnostics"
	"github.com/marcos0o0/herramienta-red/internal/repair"
	"golang.org/x/sys/windows"
)

func main() {
	if !amAdmin() {
		runMeElevated()
		return
	}
	diagnostics.CheckAdapter()
	fmt.Println("--------------------------------------------------")
	gateway := diagnostics.GetGateway()
	fmt.Println("Router detectado:", gateway)
	fmt.Println("--------------------------------------------------")
	repair.NewDHCP()
	fmt.Println("--------------------------------------------------")
	diagnostics.UsePing(gateway, "Router")
	fmt.Println("--------------------------------------------------")
	diagnostics.UsePing("8.8.8.8", "Internet")
	fmt.Println("--------------------------------------------------")
	diagnostics.UseNsLookup()
	fmt.Println("--------------------------------------------------")
	repair.ResetWinsock()
	fmt.Println("--------------------------------------------------")

}

func runMeElevated() {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 //SW_NORMAL

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		fmt.Println(err)
	}
}

func amAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	return true
}
