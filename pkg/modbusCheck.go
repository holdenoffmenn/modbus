package pkg

import (
	"fmt"
	"time"

	modbus "github.com/goburrow/modbus"
	utilsPkg "github.com/holdenoffmenn/modbus/utils"
	//mqtt "github.com/eclipse/paho.mqtt.golang"
)



func ConnModbus(device utilsPkg.DevSettings) bool {	
	var connectionTry int = 0	
	mdbsDevice := modbus.NewTCPClientHandler(device.Address + ":" + device.Port)
	for {
		err := mdbsDevice.Connect()
		if err != nil && connectionTry < 30 {
			fmt.Printf("FIRST CONNECTION: FAIL - DEVICE [%s]is Offline. Try connection[%d]  ", 
						mdbsDevice.Address, connectionTry)
			connectionTry++
			// Aguardar 5 segundos para tentar reconexão
			fmt.Println("routineMdbRead: Waiting 1 minute to try reconnect. Attempt ", connectionTry)
			time.Sleep(5 * time.Second)
		} else if connectionTry > 30 {
			fmt.Println("routineMdbRead: FAIL - Read Stop for ", mdbsDevice.Address)
			connectionTry = 0
			return false
		} else {
			return true
		}

	}

	
}