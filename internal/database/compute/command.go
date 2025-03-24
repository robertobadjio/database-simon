package compute

const (
	// SetCommand ...
	SetCommand = "SET"
	// GetCommand ...
	GetCommand = "GET"
	// DelCommand ...
	DelCommand = "DEL"
	// UnknownCommand ...
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
