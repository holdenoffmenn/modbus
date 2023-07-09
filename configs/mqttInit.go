package configs

import (
	"fmt"
	//"os"
	//utilsPkg "github.com/holdenoffmenn/modbus/utils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/holdenoffmenn/modbus/pkg"
)

var MqttMsgRsp bool
var MqttClient mqtt.Client
var MqttToken mqtt.Token
var MqttOptions *mqtt.ClientOptions

func MqttCommunication() mqtt.Client {

	if MqttBrokerInfo.Server == "" || MqttBrokerInfo.Port == "" {
		fmt.Printf("mqttInit:MqttComunication: Server[%s] or Port[%s] is empty. \n",
			MqttBrokerInfo.Server, MqttBrokerInfo.Port)
	} else {
		fmt.Printf("mqttInit:MqttCommunication: Starting communication with Broker MQTT[%s:%s]\n",
			MqttBrokerInfo.Server, MqttBrokerInfo.Port)

		var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {			
			//TODO - Recebe a mensagem e faz uma chamada. Criar as go routines com canais a partir daqui
			inMsg, err := pkg.MsgMQTTInput(msg, MqttClient)
			if err != nil {
				fmt.Println(err)
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
			if MqttToken = c.Subscribe("fails", 0, messageHandler); MqttToken.Wait() && MqttToken.Error() != nil {
				fmt.Printf("%v", MqttToken.Error())
			}
		}

		MqttClient = mqtt.NewClient(MqttOptions)
		MqttToken = MqttClient.Connect()

		if MqttToken.Wait() && MqttToken.Error() != nil {
			fmt.Printf("mqttInit:MqttCommunication: Fail to connect with Broker Ip:Port(%s:%s)\n",
				MqttBrokerInfo.Server, MqttBrokerInfo.Port)
		} else {
			fmt.Printf("mqttInit:MqttCommunication: MQTT Client is Connected to MQTT Brocker Ip:Port[%s:%s]\n",
				MqttBrokerInfo.Server, MqttBrokerInfo.Port)
		}
		//cfgPkg.SettingsMqtt = mqttSet
		return MqttClient
	}

	//cfgPkg.SettingsMqtt = mqttSet
	return MqttClient

}
