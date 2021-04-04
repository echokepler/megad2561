package megad2561

import (
	"strconv"
)

type ModeType uint8

const (
	// P - устройство реагирует (то есть отправляет сообщения на сервер, выполняет сценарии и т.д.)
	// только при замыкании контакта/выключателя
	P ModeType = iota

	// R - устройство реагирует только при размыкании контакта/выключателя.
	// На сервер отправляется дополнительный параметр "m=1".
	R

	// PR - устройство реагирует как на замыкание, так и на размыкания контакта.
	PR

	// CLICK - Click Mode (обработка одинарных и двойных кликов/нажатий)
	CLICK
)

type PortInput struct {
	*BasePort

	Commands string // @TODO need implements for command

	// ForceSendToNet eсли он не установлен (по умолчанию),
	// то сценарий выполняется ТОЛЬКО если сервер не прописан,
	// недоступен или HTTP-статус отличен от 200.
	// Если флажок установлен, то сценарий выполняется всегда независимо от наличия сервера.
	// Контроллер в этом случае будет сообщать на сервер о событиях,
	// но его ответные команды в рамках одной TCP-сессии будут проигнорированы.
	ForceSendToNet bool

	// NetCommandAddress В этом поле записывается URL, который MegaD-2561 вызывает независимо от того,
	// есть сервер или его нет. Этот URL вызывается после попытки связи с сервером и после того,
	// как отработает сценарий, описанный в поле Action. После IP-адреса можно указать порт.
	// По умолчанию 80.
	NetCommandAddress string

	// NetEnableOnlyOnFailure  указывает, что NetAction будет вызван ТОЛЬКО при недоступности сервера
	// (или когда HTTP-статус ответа отличен от 200). По умолчанию вызывается всегда.
	NetEnableOnlyOnFailure bool

	Mode ModeType

	// ModeRaw параметр отключает встроенную защиту от дребезга.
	IsRaw bool

	// ModeMute параметр отключает отправку информации на сервер о переключениях входа.
	IsMute bool
}

func (port *PortInput) Read(values ServiceValues) error {
	var err error

	if values.Has("af") {
		port.ForceSendToNet, err = strconv.ParseBool(values.Get("af"))
		if err != nil {
			return err
		}
	}

	if values.Has("naf") {
		port.NetEnableOnlyOnFailure, err = strconv.ParseBool(values.Get("naf"))
		if err != nil {
			return err
		}
	}

	if values.Has("d") {
		port.IsRaw, err = strconv.ParseBool(values.Get("d"))
		if err != nil {
			return err
		}
	}

	if values.Has("mt") {
		port.IsMute, err = strconv.ParseBool(values.Get("mt"))
		if err != nil {
			return err
		}
	}

	if values.Has("m") {
		mode, err := strconv.ParseInt(values.Get("m"), 10, 64)
		if err != nil {
			return err
		}
		port.Mode = ModeType(mode)
	}

	port.Commands = values.Get("ecmd")
	port.NetCommandAddress = values.Get("eth")

	return nil
}

func (port *PortInput) Write() (ServiceValues, error) {
	values := ServiceValues{}

	values.Add("pn", strconv.FormatInt(int64(port.BasePort.ID), 10))
	values.Add("pty", strconv.FormatInt(int64(port.GetType()), 10))
	values.Add("ecmd", port.Commands)
	values.Add("eth", port.NetCommandAddress)
	values.Add("m", strconv.FormatInt(int64(port.Mode), 10))
	values.Add("af", strconv.FormatBool(port.ForceSendToNet))
	values.Add("naf", strconv.FormatBool(port.NetEnableOnlyOnFailure))
	values.Add("d", strconv.FormatBool(port.IsRaw))
	values.Add("mt", strconv.FormatBool(port.IsMute))

	return values, nil
}
