package pkg

import (
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	utilsPkg "github.com/holdenoffmenn/modbus/utils"
)

func MsgMQTTInput(msg mqtt.Message, client mqtt.Client) (utilsPkg.MessageInput, error ){

	var msgInput utilsPkg.MessageInput
	err := json.Unmarshal(msg.Payload(), &msgInput)

	if err != nil {
		fmt.Printf("mqttInput:MsgMQTTInput: Fail to decode JSON Input Msg MQTT - Data[%s] err[%v]\n",
			msg.Payload(), err)
		return msgInput, err
	}
	fmt.Printf("mqttInput:MsgMQTTInput: msgMdbs.Protocol    [%s]\n", msgInput)	
	return msgInput, err
}
