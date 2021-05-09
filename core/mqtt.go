package core

type MegadPortInMessage struct {
	Port  int    `json:"port"`
	Mode  int8   `json:"m"`
	Value string `json:"value"`
	Click int8   `json:"click,omitempty"`
	Count int32  `json:"cnt"`
}

type MessageHandlerCallback func(msg MegadPortInMessage)

type MqttService interface {
	// DoCommand публикует сообщение в формате megad сервера
	DoCommand(commander CommandReader) error

	// SubscribePortIn подписываемся на порт
	//
	// Только для портов типа Input (InputPort)
	SubscribePortIn(ID uint8, handler MessageHandlerCallback) error
}
