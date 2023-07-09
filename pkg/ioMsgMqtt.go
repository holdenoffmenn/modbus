package pkg

import (
	"encoding/json"

	utilsPkg "github.com/holdenoffmenn/modbus/utils"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var MqttMsgRsp bool
var MqttClient mqtt.Client
var MqttToken mqtt.Token
var MqttOptions *mqtt.ClientOptions

func MQTTSendMessageAutomatic(msg utilsPkg.ExitPayloadMsg) bool {

	// data := utilsPkg.MqttMsgStruct{
	// 	Address: utilsPkg.Address{
	// 		Address: msg.Address,
	// 		Port:    msg.Port,
	// 		Name:    msg.Name,
	// 		Others:  msg.Others,
	// 	},
	// 	ReadTimeStamp: msg.ReadTimeStamp, //.Format(time.RFC3339),
	// 	Protocol:      msg.Protocol,
	// 	Data: utilsPkg.Data{
	// 		BitMemories:  msg.BitMemories,
	// 		WordMemories: msg.WordMemories,
	// 	},
	// }

	data := utilsPkg.InputMQTTMsg{
		Operation: "testando o teste",
	}

	jsonPayload, err := json.Marshal(data)
	if err != nil {
		panic(err) //Check - não pode parar o programa
	}

	if !MqttClient.IsConnected() {
		fmt.Printf("ioMsgMqtt:MQTTSendMessageAutomatic: FAIL to connect with MQTT")
	} else {
		for _, channelMqtt := range msg.Topics {
			token := MqttClient.Publish(channelMqtt, 0, false, jsonPayload)
			token.Wait()

			if token.Error() != nil {
				fmt.Printf("ioMsgMqtt:MQTTSendMessageAutomatic: FAIL to send a MQTT")
				fmt.Printf("%v", token.Error())
			}

		}
	}

	return true
}

// It sends the status of the device in a specific channel in case the connection is lost.
// func MQTTSendStatusDevice(device utilsPkg.DevSettings, settingsMqtt utilsPkg.SettingsMQTT, status bool) {

// 	type msgMqtt struct {
// 		Address  string `json:"address"`
// 		Name     string `json:"name"`
// 		Protocol string `json:"protocol"`
// 		Status   bool   `json:"status"`
// 	}

// 	msg := msgMqtt{
// 		Address:  device.Address + ":" + device.Port,
// 		Name:     device.Name,
// 		Protocol: device.Protocol,
// 		Status:   status,
// 	}

// 	jsonPayload, err := json.Marshal(msg)
// 	if err != nil {
// 		panic(err) //Check - não pode parar o programa
// 	}

// 	if !MqttClient.IsConnected() {
// 		logPkg.CtsLog.Error("ioMsgMqtt:MQTTSendStatusDevice: MQTT service is not conected!")
// 	} else {

// 		token := MqttClient.Publish(settingsMqtt.StatusDevice, 0, false, jsonPayload)
// 		token.Wait() // Aguarda a conclusão da publicação

// 		if token.Error() != nil {
// 			logPkg.CtsLog.Error("ioMsgMqtt:MQTTSendStatusDevice: Unable to send message")
// 		} else {
// 			logPkg.CtsLog.Info("ioMsgMqtt:MQTTSendStatusDevice: Sent a MQTT message.")
// 		}
// 	}
// }
