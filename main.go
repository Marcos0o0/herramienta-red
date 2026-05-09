package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
)

func main() {
	if !amAdmin() {
		runMeElevated()
		return
	}
	checkAdapter()
	fmt.Println("--------------------------------------------------")
	gateway := getGateway()
	fmt.Println("Router detectado:", gateway)
	fmt.Println("--------------------------------------------------")
	newDHCP()
	fmt.Println("--------------------------------------------------")
	usePing(gateway, "Router")
	fmt.Println("--------------------------------------------------")
	usePing("8.8.8.8", "Internet")
	fmt.Println("--------------------------------------------------")
	useNsLookup()
	fmt.Println("--------------------------------------------------")
	resetWinsock()
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

func getGateway() string {
	var ip string
	cmd := exec.Command("cmd", "/c", "chcp 65001 > nul && ipconfig /all")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	searchIpGateway := "Puerta de enlace predeterminada"
	if strings.Contains(string(output), searchIpGateway) {
		lines := strings.Split(string(output), "\n")
		for i, line := range lines {
			if strings.Contains(line, searchIpGateway) {
				nextLine := lines[i+1]
				parts := strings.Split(nextLine, ":")
				ip := strings.Trim(parts[len(parts)-1], " \r\n")
				if ip != "" && strings.Contains(ip, ".") && !strings.Contains(ip, ":") {
					return ip
				}
			}

		}
	} else {
		fmt.Println("No se ha encontrado la puerta de enlace predeterminada")
	}

	return ip
}

func usePing(gateway string, description string) bool {
	cmd := exec.Command("ping", gateway)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("No se pudo verificar la conexión:", err)
	}

	if strings.Contains(string(output), "TTL") {
		fmt.Println("Conexión a", description, "exitosa")
		return true
	} else {
		fmt.Println("Conexión a", description, "fallida")
		return false
	}
}

func useNsLookup() {
	cmd := exec.Command("nslookup", "google.com")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error al hacer conexion con google.com")
	}

	if strings.Contains(string(output), "Nombre") {
		fmt.Println("DNS funcionando correctamente")
	} else {
		fmt.Println("El DNS no está resolviendo nombres, ¿Deseas cambiarlo? (y/n)")
		userResponse := bufio.NewScanner(os.Stdin)
		userResponse.Scan()
		response := userResponse.Text()
		if response == "y" {
			cmdflush := exec.Command("ipconfig", "/flushdns").Run()
			if cmdflush != nil {
				fmt.Println("No se ha podido vaciar el DNS")
			} else {
				fmt.Println("DNS vaciado exitosamente")
			}
			adapterNameList := getAdapters("Habilitado")
			if len(adapterNameList) == 0 {
				fmt.Println("No se encontro ningun adaptador habilitado")
			} else {
				fmt.Println("Adaptadores habilitados:")
				for _, adapterName := range adapterNameList {
					fmt.Println(adapterName)
				}
				fmt.Println("¿Deseas definirle un DNS a alguno de estos adaptadores? (y/n)")
				userResponse := bufio.NewScanner(os.Stdin)
				userResponse.Scan()
				response := userResponse.Text()
				if response == "y" {
					fmt.Println("Ingrese el nombre del adaptador")
					userResponse := bufio.NewScanner(os.Stdin)
					userResponse.Scan()
					adapterName := userResponse.Text()
					cmdDNS := exec.Command("netsh", "interface", "ipv4", "set", "dnsservers", "name="+adapterName, "source=static", "address=8.8.8.8")
					if cmdDNS.Run() != nil {
						fmt.Println("No se ha podido definir el DNS")
					} else {
						fmt.Println("DNS definido exitosamente")
					}
				}
			}

		}
	}
}

func checkAdapter() {
	adapterNameList := getAdapters("Deshabilitado")
	if len(adapterNameList) > 0 {
		fmt.Println("Los adaptadores estan deshabilitados, ¿deseas activarlos? (y/n)")
		userResponse := bufio.NewScanner(os.Stdin)
		userResponse.Scan()
		response := userResponse.Text()
		if response == "y" {
			for _, name := range adapterNameList {
				enableadapter := exec.Command("netsh", "interface", "set", "interface", name, "admin=enabled").Run()
				if enableadapter != nil {
					fmt.Println("No se ha podido habilitar el adaptador", name)
				} else {
					fmt.Println("Adaptador", name, "habilitado exitosamente")
				}
			}
		} else {
			fmt.Println("No se han habilitado los adaptadores de red")
		}

	}
}

func getAdapters(status string) []string {
	cmd := exec.Command("netsh", "interface", "show", "interface")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error al obtener los adaptadores de red")
	}
	if strings.Contains(string(output), status) {
		lines := strings.Split(string(output), "\n")
		adapterNameList := []string{}
		for _, line := range lines {
			if strings.Contains(line, status) {
				fields := strings.Fields(line)
				adapterName := strings.Join(fields[3:], " ")
				adapterNameList = append(adapterNameList, adapterName)
			}
		}
		return adapterNameList
	}
	return nil
}

func newDHCP() {
	fmt.Println("¿Deseas liberar la IP actual? (y/n)")
	userResponse := bufio.NewScanner(os.Stdin)
	userResponse.Scan()
	response := userResponse.Text()
	if response == "y" {
		cmdRelease := exec.Command("ipconfig", "/release")
		if cmdRelease.Run() != nil {
			fmt.Println("Error al liberar la IP")
		} else {
			fmt.Println("IP liberada exitosamente")
			cmdRenew := exec.Command("ipconfig", "/renew")
			if cmdRenew.Run() != nil {
				fmt.Println("Error al renovar la IP")
			} else {
				fmt.Println("IP renovada exitosamente")
			}
		}
	}
}

func resetWinsock() {
	fmt.Println("¿Deseas resetear Winsock? (y/n)")
	userResponse := bufio.NewScanner(os.Stdin)
	userResponse.Scan()
	response := userResponse.Text()
	if response == "y" {
		cmdResetWinsock := exec.Command("netsh", "winsock", "reset")
		if cmdResetWinsock.Run() != nil {
			fmt.Println("Error al resetear Winsock")
		} else {
			fmt.Println("Winsock reseteado exitosamente")
			fmt.Println("Es necesario reiniciar el PC para que los cambios tomen efecto")
		}
	}
}
