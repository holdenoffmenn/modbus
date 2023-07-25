package pkg

import(
	utilsPkg "github.com/holdenoffmenn/modbus/utils"
)

func SendStatusProtocol(status string) {
	var data utilsPkg.MessageStatusProtocol
	data.MessageType = "status"
	data.Data.Status = status
	data.Data.Name = "modbus"
	Sender(data, utilsPkg.Topics)
}

func SendStatusDevice(device utilsPkg.DevSettings, status string){
	data := utilsPkg.MessageDeviceStatus{
		MessageType: "deviceStatus",
		Data: utilsPkg.MessageDataDeviceStatus{
			Protocol: "modbus",
			Device: device.Name,
			Status: status,
		},
	}
	
	Sender(data, utilsPkg.Topics)
}