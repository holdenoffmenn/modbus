package pkg

import(
	"fmt"
	utilsPkg "github.com/holdenoffmenn/modbus/utils"
	modbus "github.com/goburrow/modbus"

	//mqtt "github.com/eclipse/paho.mqtt.golang"

)



func ConnModbus(device utilsPkg.DevSettings) bool {

	//Add client to global Modbus variable
	// if MqttClient == nil {
	// 	MqttClient = client
	// }

	mdbsDevice := modbus.NewTCPClientHandler(device.Address + ":" + device.Port)

	err := mdbsDevice.Connect()
	defer mdbsDevice.Close()

	if err != nil {
		fmt.Println("FAIL to Connect with Modbus ", device.Name)
		fmt.Println(err)
		return false

	} else {
		fmt.Println("PASS - " + device.Name + " is Working in Modbus.")
		return true
	}
}