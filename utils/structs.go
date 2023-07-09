package utils

import "sync"

type ConfigMQTT struct {
	Server   string `json:"server"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type InputMQTTMsg struct {
	Operation string `json:"operation"`
	// Address   string                   `json:"address"`
	// Port      string                   `json:"port"`
	// Protocol  string                   `json:"protocol"`
	// Topics    []string                 `json:"topics"`
	// Data      []map[string]interface{} `json:"data"`
}

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
	Wg       sync.WaitGroup // WaitGroup para aguardar a conclusão das goroutines
}

// Função para encerrar as goroutines existentes
func (c *RoutineController) Stop() {
	close(c.StopChan) // Sinalize para encerrar as goroutines
	c.Wg.Wait()       // Aguarde a conclusão das goroutines
}