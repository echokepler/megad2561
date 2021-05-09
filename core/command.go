package core

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ActionType string
type TypeCommandError int

const (
	BaseAction   ActionType = "base"
	ShimAction   ActionType = "shim"
	PauseAction  ActionType = "p"
	GlobalAction ActionType = "g"
	RepeatAction ActionType = "r"
)

const (
	UnknownCommand TypeCommandError = iota
)

type CommandError struct {
	Type TypeCommandError
	Err  error
}

func (ce *CommandError) Error() string {
	return ce.Err.Error()
}

type ICommand interface {
	String() string
	GetType() ActionType
}

// ParseCommand используется при парсинге строки пришедшей от сервера megad
func ParseCommand(cmd string) (ICommand, error) {
	divided := strings.Split(cmd, ":")

	if len(divided) > 1 {
		re := regexp.MustCompile(`\D`)
		match := re.FindString(divided[1])
		target, err := strconv.ParseUint(divided[0], 10, 64)
		if err != nil {
			return nil, err
		}

		if len(match) == 0 {
			value, err := strconv.ParseUint(divided[1], 10, 64)
			if err != nil {
				return nil, err
			}

			command := BaseCommand{}
			command.Value = uint16(value)
			command.TargetPort = uint8(target)

			return &command, nil
		}

		command := ShimCommand{}
		command.TargetPort = uint8(target)
		err = command.Parse(cmd)
		if err != nil {
			return nil, err
		}

		return &command, nil
	}

	switch cmd[0] {
	case 'p':
		command := PauseCommand{}
		return &command, command.Parse(cmd)
	case 'g':
		command := GlobalCommand{}
		return &command, command.Parse(cmd)
	case 'r':
		command := RepeatCommand{}
		return &command, command.Parse(cmd)
	}

	fmt.Println(cmd)

	return nil, &CommandError{
		Type: UnknownCommand,
		Err:  errors.New("unknown command"),
	}
}

// Command  структура для удобного хранения команды
type Command struct {
	Value uint16
}

// BaseCommand базовый сценарий
// Формат поля Action следующий: X:Y;X:Y;X:Y
// где, X - номер порта, а Y - действие/команда
//
// Возможные команды:
//
// 0 - выключить;
// 1 - включить;
// 2 - изменить состояние на противоположное (переключить), т.е. если было включено выключить и наоборот
// 3 - прямая синхронизация выхода со входом (кнопка нажата - лампа включена; кнопка отпущена - лампа выключена)
// 4 - обратная синхронизация выхода со входов (кнопка нажата - лампа выключена; кнопка отпущена - лампа включена)
// [0..255] - в случае с диммируемым портами, установить значение диммера (яркости света)
type BaseCommand struct {
	Command
	TargetPort uint8
}

func (command *BaseCommand) GetType() ActionType {
	return BaseAction
}

func (command *BaseCommand) String() string {
	return fmt.Sprintf("%v:%v", command.TargetPort, command.Value)
}

type ModeShimType int

const (
	PowerToggle ModeShimType = iota
	PowerMore
	PowerLess
	PowerTransition
)

// ShimCommand Команды для управления диммируемыми каналами: +, -, ~
type ShimCommand struct {
	Command
	Mode         ModeShimType
	WithModifier bool
	TargetPort   uint8
}

func (command *ShimCommand) GetType() ActionType {
	return ShimAction
}

func (command *ShimCommand) String() string {
	format := "%v"

	switch command.Mode {
	case PowerToggle:
		format += "*%v"
	case PowerMore:
		format += "+"
	case PowerLess:
		format += "-"
	case PowerTransition:
		format += "~"
	}

	if command.WithModifier {
		format += "*"
	}

	return fmt.Sprintf(format, command.TargetPort, command.Value)
}

func (command *ShimCommand) Parse(cmd string) error {
	divided := strings.Split(cmd, ":")

	target, err := strconv.ParseUint(divided[0], 10, 64)
	if err != nil {
		return err
	}

	rePowerToggle := regexp.MustCompile(`\*[0-9]+`)
	powerToggleMatch := rePowerToggle.FindString(divided[1])

	command.TargetPort = uint8(target)

	if len(powerToggleMatch) > 0 {
		valStr := powerToggleMatch[1:]

		value, err := strconv.ParseUint(valStr, 10, 64)
		if err != nil {
			return err
		}

		command.Mode = PowerToggle
		command.Value = uint16(value)
	}

	switch divided[1][0] {
	case '+':
		command.Mode = PowerMore
	case '-':
		command.Mode = PowerLess
	case '~':
		command.Mode = PowerTransition
	}

	if len(divided[1]) > 1 && divided[1][1] == '*' {
		command.WithModifier = true
	}

	return nil
}

// PauseCommand Сценарий пауз
type PauseCommand struct {
	Command
}

func (command *PauseCommand) GetType() ActionType {
	return PauseAction
}

func (command *PauseCommand) String() string {
	return fmt.Sprintf("p%v", command.Value)
}

func (command *PauseCommand) Run() {
	time.Sleep(time.Duration(command.Value) * time.Second)
}

func (command *PauseCommand) Parse(cmd string) error {
	str := cmd[1:]
	value, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return err
	}

	command.Value = uint16(value)

	return nil
}

// GlobalCommand глобальный сценарний
type GlobalCommand struct {
	Command
}

func (command *GlobalCommand) GetType() ActionType {
	return GlobalAction
}

func (command *GlobalCommand) String() string {
	return fmt.Sprintf("g%v", command.Value)
}

func (command *GlobalCommand) Parse(cmd string) error {
	str := cmd[1:]
	value, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return err
	}

	command.Value = uint16(value)

	return nil
}

// RepeatCommand повтор сценария
type RepeatCommand struct {
	Command
}

func (command *RepeatCommand) GetType() ActionType {
	return RepeatAction
}

func (command *RepeatCommand) String() string {
	return fmt.Sprintf("r%v", command.Value)
}

func (command *RepeatCommand) Parse(cmd string) error {
	str := cmd[1:]
	value, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return err
	}

	command.Value = uint16(value)

	return nil
}
