package utils

import (
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)



type ConfigMQTT struct {
	Server   string `json:"server"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Formato da mensagem de entrada esperado pelo CallBack
// TODO: Deixar informação em repositório externo
type InputMQTTMsg struct {
	Action string `json:"action"`
}

// Criando um canal booleano

 var Wg sync.WaitGroup // WaitGroup para aguardar a conclusão das goroutines



type Devices struct {
	Devices []DevSettings `json:"devices"`
}

type DevSettings struct {
	Address     string                   `json:"address"`
	Port        string                   `json:"port"`
	Name        string                   `json:"name"`
	Protocol    string                   `json:"protocol"`
	ReadingTime int                      `json:"readingTime"`
	Topics      []string                 `json:"topics"`
	Data        []map[string]interface{} `json:"data"`
}

type ExitPayloadMsg struct {
	Address       string
	Port          string
	Name          string
	Others        string
	Protocol      string
	ReadTimeStamp string
	Topics        []string
	BitMemories   map[string]interface{}
	WordMemories  map[string]interface{}
}

var MqttClient mqtt.Client

var LoopModbus bool = true
var Status bool = true

type MqttMsgStruct struct {
	Address       Address `json:"address"`
	ReadTimeStamp string  `json:"readTimeStamp"`
	Protocol      string  `json:"protocol"`
	Data          Data    `json:"data"`
}

type Address struct {
	Address string `json:"address"`
	Port    string `json:"port"`
	Name    string `json:"name"`
	Others  string `json:"others"`
}

type Data struct {
	BitMemories  map[string]interface{} `json:"bitMemories"`
	WordMemories map[string]interface{} `json:"wordMemories"`
}

// Send Status Protocol
type StatusProtocol struct{
	Operation Info `json:"operation"`
}
type Info struct {
	Status   string `json:"status"`
	Protocol string `json:"protocol"`
}

//////////////////////////////////////////////////////////////////////////////////////////////

type MessageStatus struct{
	MessageType string `json:"messageType"`
	Data MessageDataStatus `json:"data"`
}

type MessageDataStatus struct{
	Status string `json:"status"`
	Name string `json:"name"`
}

var DoneChan chan bool

func CreateChannel(){
	DoneChan = make(chan bool)
}