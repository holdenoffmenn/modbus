package pkg

import (
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	utilsPkg "github.com/holdenoffmenn/modbus/utils"
)

func MsgMQTTInput(msg mqtt.Message, client mqtt.Client) (utilsPkg.InputMQTTMsg, error ){

	var msgInf utilsPkg.InputMQTTMsg
	err := json.Unmarshal(msg.Payload(), &msgInf)

	if err != nil {
		fmt.Printf("mqttInput:MsgMQTTInput: Fail to decode JSON Input Msg MQTT - Data[%s] err[%v]\n",
			msg.Payload(), err)
		return msgInf, err
	}
	fmt.Printf("mqttInput:MsgMQTTInput: msgMdbs.Protocol    [%s]\n", msgInf)	
	return msgInf, err
}
