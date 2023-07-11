package pkg

import (
	"time"

	utilsPkg "github.com/holdenoffmenn/modbus/utils"

	//mqtt "github.com/eclipse/paho.mqtt.golang"
	"fmt"

	modbusLib "github.com/goburrow/modbus"
	"github.com/holdenoffmenn/functions"
)

type MemoryInfoBit struct {
	Format  uint8  `json:"format"`
	Address uint16 `json:"address"`
	Name    string `json:"name"`
	Value   bool   `json:"value"`
}

type MemoryInfoWord struct {
	Format  uint8  `json:"format"`
	Address uint16 `json:"address"`
	Name    string `json:"name"`
	Value   uint32 `json:"value"`
}
type PlcInfo struct {
	BitMemories  []MemoryInfoBit  `json:"bitMemories"`
	SlaveID      uint8            `json:"slaveId"`
	WordMemories []MemoryInfoWord `json:"wordMemories"`
}

var connectionTry int = 0
var readHour time.Time

// ReadInfoMdbs will read informations from device
func ReadInfoMdbs(dev utilsPkg.DevSettings) {
	//Start Connection
	conn := modbusLib.NewTCPClientHandler(dev.Address + ":" + dev.Port)
	conn.Timeout = 3 * time.Minute

	err := conn.Connect()
	defer conn.Close()

	plcInfo := changeFormatRead(dev)

	//Armazena resultados
	resultsGetBitMemories := make(map[string]interface{})
	resultsGetWordMemories := make(map[string]interface{})
	lastBitMemories := make(map[string]interface{})
	lastWordMemories := make(map[string]interface{})

	if err != nil {
		fmt.Printf("routineMdbRead: FAIL to Connect with Modbus device [%s] error [%s]", dev.Name, err)
		return
	} else {

		clientModbus := modbusLib.NewClient(conn)
		readHour = time.Now()
		readTime := readHour.Format("2006-01-02 15:04:05")

		//loop infinito de verificação, repetindo conforme o tempo informado no arquivo config.json
		for utilsPkg.LoopModbus {
			fmt.Println("-----> AINDA RODANDO ==> ", dev.Name)
			fmt.Printf("routineMdbRead: dev name[%s]  Address[%s:%s] readTime[%v]Reading\n",
				dev.Name, dev.Address, dev.Port, readTime)
			if len(plcInfo.BitMemories) != 0 {

				// Varre items configurados no arquivo config.json
				for _, itemMemories := range plcInfo.BitMemories {
					resultBit, err := readBitMemory(uint8(plcInfo.SlaveID), clientModbus, itemMemories.Address)
					if err != nil {
						//Forçar o fechamento da conexão
						conn.Close()
						//Adiciona FAIL como valor na chave que houve falha de leitura
						resultsGetBitMemories[itemMemories.Name] = resultBit
						//Chama o loop de verificação de conexão
						_ = connCheck(conn)
					} else {
						// Salva o resultado encontrado
						resultsGetBitMemories[itemMemories.Name] = resultBit
					}

				}

			}
			if len(plcInfo.WordMemories) != 0 {
				// Varre items configurados no arquivo config.json
				for _, itemMemories := range plcInfo.WordMemories {
					resultWord, err := readWordMemory(
						plcInfo.SlaveID,
						clientModbus,
						itemMemories.Address,
						itemMemories.Format)

					if err != nil {
						//Forçar o fechamento e reabertura da conexão
						conn.Close()
						resultsGetWordMemories[itemMemories.Name] = resultWord
						_ = connCheck(conn)
					} else {
						resultsGetWordMemories[itemMemories.Name] = resultWord
					}
				}

			}

			changesBit := functions.CompareMapsStrIFace(lastBitMemories, resultsGetBitMemories)
			changesWord := functions.CompareMapsStrIFace(lastWordMemories, resultsGetWordMemories)

			if !changesBit && !changesWord {
				// Não houve alterações
				fmt.Printf("routineMdbRead: dev name[%s]  Ip[%s:%s] No  Changes Oddurred - Sleeping %v Seconds...\n",
					dev.Name, dev.Address, dev.Port, dev.ReadingTime)
			} else {
				// Holve alguma alteração gerar payload JSON da resposta
				msg := utilsPkg.ExitPayloadMsg{
					Address:       dev.Address,
					Port:          dev.Port,
					Name:          dev.Name,
					Others:        "",
					Protocol:      dev.Protocol,
					ReadTimeStamp: readTime,
					Topics:        dev.Topics,
					BitMemories:   resultsGetBitMemories,
					WordMemories:  resultsGetWordMemories,
				}
				fmt.Printf("Bit results in [%s] - [%v]\n", dev.Name, resultsGetBitMemories)
				fmt.Printf("Word results in [%s] - [%v]\n", dev.Name, resultsGetWordMemories)

				// Enviar payload json da mensagem
				_ = MQTTSendMessageAutomatic(msg)
			}

			if !utilsPkg.LoopModbus {
				// Thread do dispositivo foi encerrada
				fmt.Printf("routineMdbRead: dev name[%s]  Ip[%s:%s] Encerrada em %v\n",
					dev.Name, dev.Address, dev.Port, dev.ReadingTime)
				return
			}

			// Atualiza ultimos valores com os valores correntes
			lastBitMemories = functions.CopyMap(resultsGetBitMemories).(map[string]interface{})
			lastWordMemories = functions.CopyMap(resultsGetWordMemories).(map[string]interface{})

			// Dorme conforme tempo configurado no arquivo json.config
			time.Sleep(time.Duration(dev.ReadingTime) * time.Second)

			readHour = time.Now()
			readTime = readHour.Format("2006-01-02 15:04:05")
		}

		fmt.Println("SAIU DO LOOP")

		//}

	}
}

// connCheck will verify if the device is running in network
// func connCheck(conn *modbusLib.TCPClientHandler, plcIP string, plcPort string) bool {
func connCheck(conn *modbusLib.TCPClientHandler) bool {
	for {
		err := conn.Connect()
		if err != nil && connectionTry < 30 {
			fmt.Println("routineMdbRead: FAIL - PLC lost connection ", conn.Address)
			connectionTry++
			// Aguardar 5 segundos para tentar reconexão
			fmt.Println("routineMdbRead: Waiting 1 minute to try reconnect. Attempt ", connectionTry)
			time.Sleep(5 * time.Second)
		} else if connectionTry > 30 {
			fmt.Println("routineMdbRead: FAIL - Read Stop for ", conn.Address)
			connectionTry = 0
			return false
		} else {
			return true
		}

	}
}

// readBitMemory will make a read to device modbus by Bit
func readBitMemory(slaveID uint8, client modbusLib.Client, addressMemory uint16) (bool, error) {
	results, err := client.ReadCoils(addressMemory, 1)
	if err != nil {
		fmt.Printf("routineMdbRead: client.ReadCoils(slaveid[%v], memoryAddress[%v], err[%s]\n", slaveID, addressMemory, err)
		return false, err
	}

	//Convert a resposta de array de byte > int > string
	if results[0] == 1 {
		return true, nil
	}
	return false, nil
}

// readWordMemory will make a read to device modbus by Word
func readWordMemory(slaveID uint8, client modbusLib.Client, addressMemory uint16, typeMem uint8) (uint32, error) {
	var resp uint32
	results, err := client.ReadHoldingRegisters(addressMemory, 1)
	if err != nil {
		fmt.Printf("%v", err)
		return resp, err
	}

	if len(results) == 4 {
		value32 := uint32(results[0])<<24 | uint32(results[1])<<16 | uint32(results[2])<<8 | uint32(results[3])
		//CtsLog.Info("Read 32-bit value: \n", value32)
		//resp = strconv.Itoa(int(value32))
		resp = value32
	} else if len(results) == 2 {
		value16 := uint32(results[0])<<8 | uint32(results[1])
		//result1 := uint32(results[0])
		//CtsLog.Info("Read 16-bit value: \n", value16)
		//resp = strconv.Itoa(int(value16))
		resp = value16
	}

	return resp, nil
}

func changeFormatRead(dev utilsPkg.DevSettings) PlcInfo {

	var bitMemories []MemoryInfoBit
	var wordMemories []MemoryInfoWord
	var slave float64

	for _, data := range dev.Data {
		bit, _ := data["bitMemories"].([]interface{})
		word, _ := data["wordMemories"].([]interface{})
		slave, _ = data["slaveId"].(float64)

		for _, item := range bit {
			dataMap, _ := item.(map[string]interface{})
			address, _ := dataMap["address"].(float64)
			name, _ := dataMap["name"].(string)

			memory := MemoryInfoBit{
				Address: uint16(address),
				Name:    name,
			}

			bitMemories = append(bitMemories, memory)
		}

		for _, item := range word {
			dataMap, _ := item.(map[string]interface{})
			address, _ := dataMap["address"].(float64)
			name, _ := dataMap["name"].(string)
			//typeMem, _ := dataMap["type"].(float64)

			memory := MemoryInfoWord{
				Address: uint16(address),
				Name:    name,
				//Format:  typeMem,
			}

			wordMemories = append(wordMemories, memory)
		}

	}

	plcInfo := PlcInfo{
		BitMemories:  bitMemories,
		SlaveID:      uint8(slave),
		WordMemories: wordMemories,
	}

	return plcInfo
}
