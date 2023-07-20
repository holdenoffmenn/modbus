package configs

import (
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/holdenoffmenn/modbus/pkg"
	utilsPkg "github.com/holdenoffmenn/modbus/utils"
)

var MqttMsgRsp bool
var MqttToken mqtt.Token
var MqttOptions *mqtt.ClientOptions

var MqttBrokerInfo utilsPkg.ConfigMQTT
var Hostname string

func SetMqttBroker() {
	fmt.Println("Starting assign values to public variables.")
	MqttBrokerInfo = utilsPkg.ConfigMQTT{
		Server:   utilsPkg.MqttAddress,
		Port:     utilsPkg.MqttPort,
		Username: utilsPkg.MqttUser,
		Password: utilsPkg.MqttPassword,
	}

	Hostname, _ = os.Hostname()
}

func MqttCommunication() error {

	if MqttBrokerInfo.Server == "" || MqttBrokerInfo.Port == "" {
		fmt.Printf("mqttInit:MqttComunication: Server[%s] or Port[%s] is empty. \n",
			MqttBrokerInfo.Server, MqttBrokerInfo.Port)
	} else {
		fmt.Printf("mqttInit:MqttCommunication: Starting communication with Broker MQTT[%s:%s]\n",
			MqttBrokerInfo.Server, MqttBrokerInfo.Port)

		var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
			inMsg, err := pkg.MsgMQTTInput(msg, utilsPkg.MqttClient)
			if err != nil {
				fmt.Printf("mqttInit:MqttCommunication: Receiving MQTT message failed. Device:[%s:%s] Error:[%s]\n",
					MqttBrokerInfo.Server, MqttBrokerInfo.Port, err)
			} else {
				DecodeMsg(inMsg)
			}
		}

		MqttOptions = mqtt.NewClientOptions()
		MqttOptions.AddBroker(MqttBrokerInfo.Server + ":" + MqttBrokerInfo.Port)
		MqttOptions.SetClientID(Hostname + uuid.New().String())
		MqttOptions.SetUsername(MqttBrokerInfo.Username)
		MqttOptions.SetPassword(MqttBrokerInfo.Password)

		MqttOptions.OnConnect = func(c mqtt.Client) {
			//Subscriber in error and write channels
			if MqttToken = c.Subscribe("modbus", 0, messageHandler); MqttToken.Wait() && MqttToken.Error() != nil {
				fmt.Printf("%v", MqttToken.Error())
			}
		}

		utilsPkg.MqttClient = mqtt.NewClient(MqttOptions)
		MqttToken = utilsPkg.MqttClient.Connect()

		if MqttToken.Wait() && MqttToken.Error() != nil {
			fmt.Printf("mqttInit:MqttCommunication: Fail to connect with Broker Ip:Port(%s:%s)\n",
				MqttBrokerInfo.Server, MqttBrokerInfo.Port)
			return MqttToken.Error()
			//TODO: Verificar como ficar tentando a conexão por um tempo
		} else {
			fmt.Printf("mqttInit:MqttCommunication: MQTT Client is Connected to MQTT Brocker Ip:Port[%s:%s]\n",
				MqttBrokerInfo.Server, MqttBrokerInfo.Port)
		}

		return nil
	}
	return fmt.Errorf("Missing_Informations")

}

func DecodeMsg(msg utilsPkg.InputMQTTMsg) {
	switch msg.Action {
	case "start":
		fmt.Println("Start Modbus")
		StartModbus()
	case "restart":
		//utilsPkg.DoneChan <- false
		fmt.Println("Reboot Modbus")
		utilsPkg.LoopModbus = false
		utilsPkg.Wg.Wait()
		StartModbus()
	case "stop":
		fmt.Println("Stop Modbus")
		utilsPkg.LoopModbus = false
		utilsPkg.Wg.Wait()
		//close(utilsPkg.StopChan)
		// utilsPkg.WgGoroutines.Wait()

		utilsPkg.LoopModbus = false
		fmt.Println("Leitura do Modbus concluída.")

	default:
		fmt.Println("Wrong Command")
	}
}
