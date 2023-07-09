package configs

import (
	//"encoding/json"
	//"fmt"
	"fmt"
	"os"
	//"sync"

	utilsPkg "github.com/holdenoffmenn/modbus/utils"
)

// type RoutineController struct {
// 	stopChan chan struct{} // Canal de sinalização para encerrar as goroutines
// 	wg       sync.WaitGroup // WaitGroup para aguardar a conclusão das goroutines
// }

//var pathConfigFileMQTT string = "configMqtt.json"

var MqttBrokerInfo utilsPkg.ConfigMQTT
var Hostname string

func Starting() {
	SetVar()
	MqttCommunication()

	

}

// SetVar will assign values to public variables
func SetVar() {
	fmt.Println("Starting assign values to public variables.")
	MqttBrokerInfo = utilsPkg.ConfigMQTT{
		Server:   "localhost",
		Port:     "1883",
		Username: "",
		Password: "",
	}

	Hostname, _ = os.Hostname()

}

func DecodeMsg(msg utilsPkg.InputMQTTMsg) {
	controller := &utilsPkg.RoutineController{
		StopChan: make(chan struct{}),
	}

	switch msg.Operation {
	case "start":
		fmt.Println("Start Modbus")
		StartModbus(controller)
	case "restart":
		fmt.Println("Reboot Modbus")
		restartRoutines(controller)
	case "stop":
		fmt.Println("Stop Modbus")
		utilsPkg.Loop = false
	default:
		fmt.Println("Wrong Command")
	}

}



// Função para reiniciar as goroutines com base no arquivo atual
func restartRoutines(controller *utilsPkg.RoutineController) {
	// Encerre as goroutines existentes
	controller.Stop()

	// Inicie as novas goroutines com base no arquivo atual
	StartModbus(controller)
	fmt.Println("Goroutines reiniciadas com base no novo arquivo.")
}


