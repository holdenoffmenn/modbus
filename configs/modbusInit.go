package configs

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	packages "github.com/holdenoffmenn/modbus/pkg"
	utilsPkg "github.com/holdenoffmenn/modbus/utils"
)

// "encoding/json"
// "fmt"
// "os"

//var pathConfigFileMQTT string
var devInf []utilsPkg.Devices
var itemInf []utilsPkg.DevSettings

func StartModbus(controller *utilsPkg.RoutineController) {
	//Ler as configurações do arquivo
	devices := GetDevConfig()
	StartRead(devices, controller)

	// Inicie as goroutines iniciais com base no arquivo atual
	//startInitialRoutines("dados.txt", controller)

}

func GetDevConfig() []utilsPkg.DevSettings{

	var devConfig string = "plcConfig.json"
	file, err := os.ReadFile(devConfig)
	if err != nil {
		fmt.Println("Fail to read JSON File: ", err)
		return nil
	}

	err = json.Unmarshal(file, &devInf)
	if err != nil {
		fmt.Println("FAIL to decode JSON file")
		fmt.Printf("%v", err)
		return nil
	}

	for _, dev := range devInf {
		for _, item := range dev.Devices {
			itemInf = append(itemInf,
				utilsPkg.DevSettings{
					Address:     item.Address,
					Port:        item.Port,
					Name:        item.Name,
					Protocol:    item.Protocol,
					ReadingTime: item.ReadingTime,
					Topics:      item.Topics,
					Data:        item.Data,
				})
		}
	}
	return itemInf
}


func StartRead(devices []utilsPkg.DevSettings, controller *utilsPkg.RoutineController) {
	for _, device := range devices {
		if device.Protocol == "modbus" {
			fmt.Printf("StartRead: Name[%s] Protocol[%s] Ip:Port[%s:%s]\n",
				device.Name, device.Protocol, device.Address, device.Port)
			
			statusPlc := packages.ConnModbus(device)
			if statusPlc {
				//modbusPkg.MQTTSendStatusDevice(device, SettingsMqtt, statusPlc)
				go packages.ReadInfoMdbs(device, controller)
			} else {
				fmt.Println("Erro - Dispositivo não está online")
				//modbusPkg.MQTTSendStatusDevice(device, SettingsMqtt, statusPlc)
			}

		}
		
		time.Sleep(1 * time.Second)
	}

}
