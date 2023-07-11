package main

import (
	"fmt"
	"time"

	configsPkg "github.com/holdenoffmenn/modbus/configs"
)

func main() {
	fmt.Println("Starting MODBUS...")
	err := startingLocalBroker()
	if err != nil {
		fmt.Println("FAIL - Is not possible start communication with broker")
		return
	}
	configsPkg.StartModbus()

	for {
		fmt.Println("Checking...")
		time.Sleep(20000 * time.Millisecond)
	}
}

func startingLocalBroker() error {
	configsPkg.SetMqttBroker()
	err := configsPkg.MqttCommunication()
	return err
}
