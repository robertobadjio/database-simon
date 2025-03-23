package compute

const (
	SetCommand     = "SET"
	GetCommand     = "GET"
	DelCommand     = "DEL"
	UnknownCommand = "UNKNOWN"
)

const (
	setCommandArgumentsNumber = 2
	getCommandArgumentsNumber = 1
	delCommandArgumentsNumber = 1
)

var argumentsNumber = map[string]int{
	SetCommand: setCommandArgumentsNumber,
	GetCommand: getCommandArgumentsNumber,
	DelCommand: delCommandArgumentsNumber,
}

func getCommand(command string) string {
	switch command {
	case SetCommand:
		return SetCommand
	case GetCommand:
		return GetCommand
	case DelCommand:
		return DelCommand
	default:
		return UnknownCommand
	}
}

func commandArgumentsNumber(command string) int {
	return argumentsNumber[command]
}
