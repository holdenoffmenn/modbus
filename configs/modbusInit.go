package configs

import (
	"encoding/json"
	"fmt"
	"os"

	//"sync"
	"time"
	"github.com/holdenoffmenn/modbus/pkg"
	packages "github.com/holdenoffmenn/modbus/pkg"
	utilsPkg "github.com/holdenoffmenn/modbus/utils"
)

var devInf []utilsPkg.Devices

func StartModbus() {
	utilsPkg.CreateChannel()
	//Read a file from root
	devices, err := GetDevConfig()
	if err != nil {
		fmt.Printf("FAIL - Unable to capture data from deviceConfig.json file. Error [%s]", err)
		return
	}
	//Start a read devices
	StartRead(devices)

}

func GetDevConfig() ([]utilsPkg.DevSettings, error) {
	var itemInf []utilsPkg.DevSettings
	file, err := os.ReadFile(utilsPkg.FilePath)
	if err != nil {
		fmt.Println("Fail to read JSON File: ", err)
		return nil, err
	}

	err = json.Unmarshal(file, &devInf)
	if err != nil {
		fmt.Println("FAIL to decode JSON file")
		fmt.Printf("%v", err)
		return nil, err
	}

	//Seleciona apenas os dispositivos com protocolo
	for _, dev := range devInf {
		for _, item := range dev.Devices {
			if item.Protocol == "modbus" {
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
	}
	return itemInf, nil
}

func StartRead(devices []utilsPkg.DevSettings) {
	for _, device := range devices {
		fmt.Printf("StartRead: Name[%s] Protocol[%s] Ip:Port[%s:%s]\n",
			device.Name, device.Protocol, device.Address, device.Port)

		statusPlc := packages.ConnModbus(device)

		if statusPlc {			
			pkg.SendStatusDevice(device, "connected")
			utilsPkg.Wg.Add(1)
			identifier := rune(device.Name[0])
			go func(device utilsPkg.DevSettings) {
				packages.ReadInfoMdbs(device, identifier)
			}(device)

		} else {
			fmt.Println("Erro - Dispositivo não está online")
			pkg.SendStatusDevice(device, "disconnected")
		}

		time.Sleep(1 * time.Second)
	}

}

