package core

import (
	"strings"
)

type RemoteCommander struct {
	source   string
	commands []ICommand
}

func NewRemoteCommander(source string) (*RemoteCommander, error) {
	commander := RemoteCommander{source: source}

	return &commander, commander.Parse()
}

// Convert форматирует команды в строку для сервиса
func (rc *RemoteCommander) Convert() string {
	return rc.source
}

// Parse преобразует из source набор команд
func (rc *RemoteCommander) Parse() error {
	commands := strings.Split(rc.source, ";")

	for _, cmd := range commands {
		command, err := ParseCommand(cmd)
		if err != nil {
			ce, ok := err.(*CommandError)

			if ok && ce.Type == UnknownCommand {
				continue
			}

			return err
		}

		rc.commands = append(rc.commands, command)
	}

	return nil
}

// Add базовая команда которая генерирует действие на смену состояния.
//
// Значения command:
//
// 0 - выключить
//
// 1 - включить
//
// 2 - изменить состояние на противоположное (переключить), т.е. если было включено выключить и наоборот.
//
// 3 - прямая синхронизация выхода со входом (кнопка нажата - лампа включена; кнопка отпущена - лампа выключена).
//
// 4 - обратная синхронизация выхода со входов (кнопка нажата - лампа выключена; кнопка отпущена - лампа включена)
//
// [0..255] - в случае с диммируемым портами, установить значение диммера (яркости света).
func (rc *RemoteCommander) Add(command BaseCommand) {
	rc.addCommand(&command)
}

func (rc *RemoteCommander) AddShim(command ShimCommand) {
	rc.addCommand(&command)
}

// AddPause В сценариях контроллер поддерживает работу с паузами.
//
// Value 1 - 100 милисекунд
//
// Паузы в полном объеме и без ограничений работают только в сценариях по умолчанию (Action).
// Начиная с версии прошивки 4.16b8 паузы также поддерживаются и в командах, поступающих извне.
// Но в этом случае одновременно может выполняться только один сценарий, содержащий паузы.
//
// При выполнении сценария, содержащего паузу, работа контроллера не блокируется. Паузы выполняются в фоновом режиме.
func (rc *RemoteCommander) AddPause(command PauseCommand) {
	rc.addCommand(&command)
}

// AddRepeat Повтор сценария.
//
//
// Повтор записанного сценария несколько раз.
//
// Включить порт 7; пауза 0,5с; выключить порт 7; пауза 0,5с; повторить все это с самого начала еще 4 раза
// Таким образом порт включится и выключится 5 раз.
// Такую команду можно использовать для более компактной записи сложных сценариев.
func (rc *RemoteCommander) AddRepeat(command RepeatCommand) {
	rc.addCommand(&command)
}

// AddGlobal Управление всеми выходами.
//
// Например: Value = 0 (выключить все выходы), Value = 1 (включить все выходы).
func (rc *RemoteCommander) AddGlobal(command GlobalCommand) {
	rc.addCommand(&command)
}

// Private

func (rc *RemoteCommander) addCommand(command ICommand) {
	rc.commands = append(rc.commands, command)
}
