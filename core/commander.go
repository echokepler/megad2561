package core

type CommandReader interface {
	Convert() string
	Parse() error

	Add(command BaseCommand)
	AddShim(command ShimCommand)
	AddPause(command PauseCommand)
	AddRepeat(command RepeatCommand)
	AddGlobal(command GlobalCommand)
}

type ExecuteCommander interface {
	Execute(cmd CommandReader)
}
