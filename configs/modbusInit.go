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
	devices, err := GetDevConfig()
	if err != nil {
		fmt.Printf("FAIL - Unable to capture data from deviceConfig.json file. Error [%s]", err)
		return
	}
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
			utilsPkg.Wg.Add(1)
			go func(device utilsPkg.DevSettings) {
				packages.ReadInfoMdbs(device)
			}(device)
			
		} else {
			fmt.Println("Erro - Dispositivo não está online")
			//modbusPkg.MQTTSendStatusDevice(device, SettingsMqtt, statusPlc)
		}

		time.Sleep(1 * time.Second)
	}

}

func SendStatusProtocol(status string) {
	var data utilsPkg.MessageStatus
	data.MessageType = "status"
	data.Data.Status = status
	data.Data.Name = "modbus"
	pkg.Sender(data, utilsPkg.Topics)
}
