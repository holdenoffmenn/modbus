package main

import (
	"fmt"
	"time"

	configsPkg "github.com/holdenoffmenn/modbus/configs"
	"github.com/holdenoffmenn/modbus/pkg"
)

func main() {
	fmt.Println("Starting MODBUS...")
	err := startingLocalBroker()
	if err != nil {
		fmt.Println("FAIL - Is not possible start communication with broker")
		return
	}
	pkg.SendStatusProtocol("running")
	time.Sleep(500 * time.Millisecond)
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
