package repair

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func NewDHCP() {
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

func ResetWinsock() {
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
