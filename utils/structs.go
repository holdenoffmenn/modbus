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

var StatusProtocol bool = true

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
// type StatusProtocol struct{
// 	Operation Info `json:"operation"`
// }
// type Info struct {
// 	Status   string `json:"status"`
// 	Protocol string `json:"protocol"`
// }

// Decode da mensagem de entrada
type MessageInput struct {
	MessageType string      `json:"messageType"`
	Data        interface{} `json:"data"`
}

// End input message

// Formato da mensagem de status do protocolo
type MessageStatusProtocol struct {
	MessageType string            `json:"messageType"`
	Data        MessageDataStatus `json:"data"`
}

type MessageDataStatus struct {
	Status string `json:"status"`
	Name   string `json:"name"`
}

//Fim

// Formato da mensagem de status do device
type MessageDeviceStatus struct {
	MessageType string                  `json:"messageType"`
	Data        MessageDataDeviceStatus `json:"data"`
}

type MessageDataDeviceStatus struct {
	Protocol string `json:"protocol"`
	Device   string `json:"device"`
	Status   string `json:"status"`
}

//Fim

// Funções do channel de controle das Go Routines
var (
	IdentifierRoutine rune
	DoneChan chan bool
	// ChanOnce sync.Once
	Wg       sync.WaitGroup
	StopChannels = make(map[rune]chan struct{})
	
)

func CreateChannel() {
	DoneChan = make(chan bool)
	StopChannels[IdentifierRoutine] = make(chan struct{})
}

//End Go Routines Control
