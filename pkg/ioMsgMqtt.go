package pkg

import (
	"encoding/json"

	"fmt"

	utilsPkg "github.com/holdenoffmenn/modbus/utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var MqttMsgRsp bool
var MqttToken mqtt.Token
var MqttOptions *mqtt.ClientOptions

func MQTTSendMessageAutomatic(msg utilsPkg.ExitPayloadMsg) bool {

	data := utilsPkg.MqttMsgStruct{
		Address: utilsPkg.Address{
			Address: msg.Address,
			Port:    msg.Port,
			Name:    msg.Name,
			Others:  msg.Others,
		},
		ReadTimeStamp: msg.ReadTimeStamp, //.Format(time.RFC3339),
		Protocol:      msg.Protocol,
		Data: utilsPkg.Data{
			BitMemories:  msg.BitMemories,
			WordMemories: msg.WordMemories,
		},
	}

	// data := utilsPkg.InputMQTTMsg{
	// 	Operation: "testando o teste",
	// }

	jsonPayload, err := json.Marshal(data)
	if err != nil {
		panic(err) //Check - não pode parar o programa
	}

	if !utilsPkg.MqttClient.IsConnected() {
		fmt.Printf("ioMsgMqtt:MQTTSendMessageAutomatic: FAIL to connect with MQTT")
	} else {
		for _, channelMqtt := range msg.Topics {
			token := utilsPkg.MqttClient.Publish(channelMqtt, 0, false, jsonPayload)
			token.Wait()

			if token.Error() != nil {
				fmt.Printf("ioMsgMqtt:MQTTSendMessageAutomatic: FAIL to send a MQTT")
				fmt.Printf("%v", token.Error())
			}

		}
	}

	return true
}



func StatusMessageSender(payload utilsPkg.MessageStatusProtocol) {

}

func Sender(data interface{}, topics []string) bool {

	jsonPayload, err := json.Marshal(data)
	if err != nil {
		panic(err) //Check - não pode parar o programa
	}

	if !utilsPkg.MqttClient.IsConnected() {
		fmt.Printf("ioMsgMqtt:MQTTSendMessageAutomatic: FAIL to connect with MQTT")
	} else {
		for _, channelMqtt := range topics {
			token := utilsPkg.MqttClient.Publish(channelMqtt, 0, false, jsonPayload)
			token.Wait()

			if token.Error() != nil {
				fmt.Printf("ioMsgMqtt:MQTTSendMessageAutomatic: FAIL to send a MQTT")
				fmt.Printf("%v", token.Error())
			}

		}
	}

	return true
}
