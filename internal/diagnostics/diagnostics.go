package diagnostics

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func GetGateway() string {
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

func UsePing(gateway string, description string) bool {
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

func UseNsLookup() {
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
			adapterNameList := GetAdapters("Habilitado")
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

func CheckAdapter() {
	adapterNameList := GetAdapters("Deshabilitado")
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

func GetAdapters(status string) []string {
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
