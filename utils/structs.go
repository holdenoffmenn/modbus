package utils

import (
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Informações do MQTT Broker Local
var MqttAddress string = "localhost"
var MqttPort string = "1883"
var MqttUser string = ""
var MqttPassword string = ""

type ConfigMQTT struct {
	Server   string `json:"server"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//Formato da mensagem de entrada esperado pelo CallBack
//TODO: Deixar informação em repositório externo
type InputMQTTMsg struct {
	Action string `json:"action"`
}

	var StopChan chan struct{} // Canal de sinalização para encerrar as goroutines
	var DoneChan chan bool
	var Wg       sync.WaitGroup // WaitGroup para aguardar a conclusão das goroutines
	var WgGoroutines sync.WaitGroup


var Loop bool = true

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

type RoutineController struct {
	StopChan chan struct{} // Canal de sinalização para encerrar as goroutines
	DoneChan chan struct{}
	Wg       sync.WaitGroup // WaitGroup para aguardar a conclusão das goroutines
}

// Função para encerrar as goroutines existentes
func (c *RoutineController) Stop() {
	close(c.StopChan) // Sinalize para encerrar as goroutines
	//c.Wg.Wait()       // Aguarde a conclusão das goroutines
	<-c.DoneChan
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
