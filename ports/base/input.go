package base

import (
	"github.com/echokepler/megad2561/core"
	"github.com/echokepler/megad2561/internal/qsparser"
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

type InputSettings struct {
	Commands string `qs:"ecmd"` // @TODO need implements for command

	// ForceSendToNet eсли он не установлен (по умолчанию),
	// то сценарий выполняется ТОЛЬКО если сервер не прописан,
	// недоступен или HTTP-статус отличен от 200.
	// Если флажок установлен, то сценарий выполняется всегда независимо от наличия сервера.
	// Контроллер в этом случае будет сообщать на сервер о событиях,
	// но его ответные команды в рамках одной TCP-сессии будут проигнорированы.
	ForceSendToNet bool `qs:"af"`

	// NetCommandAddress В этом поле записывается URL, который MegaD-2561 вызывает независимо от того,
	// есть сервер или его нет. Этот URL вызывается после попытки связи с сервером и после того,
	// как отработает сценарий, описанный в поле Action. После IP-адреса можно указать порт.
	// По умолчанию 80.
	NetCommandAddress string `qs:"eth"`

	// NetEnableOnlyOnFailure  указывает, что NetAction будет вызван ТОЛЬКО при недоступности сервера
	// (или когда HTTP-статус ответа отличен от 200). По умолчанию вызывается всегда.
	NetEnableOnlyOnFailure bool `qs:"naf"`

	Mode ModeType `qs:"m"`

	// ModeRaw параметр отключает встроенную защиту от дребезга.
	IsRaw bool `qs:"d"`

	// ModeMute параметр отключает отправку информации на сервер о переключениях входа.
	IsMute bool `qs:"mt"`
}

type InputPort struct {
	*Port
	settings InputSettings
}

func (port *InputPort) read(values core.ServiceValues) error {
	return qsparser.UnMarshal(values, &port.settings)
}

func (port *InputPort) write() (core.ServiceValues, error) {
	values := qsparser.Marshal(port.settings)

	return values, nil
}
