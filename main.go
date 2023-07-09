package main

import (
	"fmt"
	"time"
	utilsPkg "github.com/holdenoffmenn/modbus/utils"
	configsPkg "github.com/holdenoffmenn/modbus/configs"
)

func main() {
	
	fmt.Println("Starting MODBUS...")
	//conectar no broker
	configsPkg.Starting()

	for utilsPkg.Loop {
		fmt.Println("Checking...")
		time.Sleep(1000 * time.Millisecond)
	}
}
