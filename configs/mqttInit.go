package configs

import (
	"fmt"
	"os"
	"time"

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
			if MqttToken = c.Subscribe(MqttOptions.ClientID+"general", 0, messageHandler); MqttToken.Wait() && MqttToken.Error() != nil {
				fmt.Printf("%v", MqttToken.Error())
			}
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

func RunAction(msg utilsPkg.MessageInput) {
	action, err := msg.Data.(map[string]interface{})["action"].(string)
	if !err {
		fmt.Println("Campo 'action' não é uma string ou está nomeado errado")
		return
	}

	protocol, err := msg.Data.(map[string]interface{})["protocol"].(string)
	if !err {
		fmt.Println("Campo 'protocol' não é uma string ou está nomeado errado")
		return
	}

	if protocol == "modbus" {
		switch action {
		//Start the read if the s
		case "start":
			if !utilsPkg.StatusProtocol {
				fmt.Println("Start Modbus")
				StartModbus()
			} else {
				fmt.Println("Rotinas já estão iniciadas")
			}

		case "restart":
			if utilsPkg.StatusProtocol {
				close(utilsPkg.DoneChan)
				utilsPkg.Wg.Wait()
				fmt.Println("Reboot Modbus")
				time.Sleep(1000 * time.Millisecond)
				StartModbus()
			} else {
				fmt.Println("Reboot Modbus")
				time.Sleep(1000 * time.Millisecond)
				StartModbus()
			}
			//utilsPkg.StatusProtocol = true
		case "stop":
			if utilsPkg.StatusProtocol {
				fmt.Println("Stop Modbus")
				close(utilsPkg.DoneChan)
				utilsPkg.Wg.Wait()
				fmt.Println("Leitura do Modbus concluída.")

			} else {
				fmt.Println("Canal já está fechado")

			}
		case "123":

			input := "A"
			identifier := rune(input[0])
			if ch, ok := utilsPkg.StopChannels[identifier]; ok {
				close(ch)
				utilsPkg.Wg.Wait()
				delete(utilsPkg.StopChannels, identifier)
				fmt.Printf("Goroutine %s stopped.\n", "input")
			}
			

		default:
			fmt.Println("Wrong Command")
		}
	}
}

func DecodeMsg(msg utilsPkg.MessageInput) {
	switch msg.MessageType {
	case "action":
		fmt.Println("Processando a mensagem recebida!")
		RunAction(msg)
	default:
		fmt.Println("Comando desconhecido!")
	}
}
