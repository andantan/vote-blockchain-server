package util

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
)

// Fatal
func RedString(s string) string {
	return Red + s + Reset
}

// Commit
func GreenString(s string) string {
	return Green + s + Reset
}

// Block
func YellowString(s string) string {
	return Yellow + s + Reset
}

func BlueString(s string) string {
	return Blue + s + Reset
}

// Init, System
func MagentaString(s string) string {
	return Magenta + s + Reset
}

// Pending
func CyanString(s string) string {
	return Cyan + s + Reset
}

func WhiteString(s string) string {
	return White + s + Reset
}

func InitString(s string) string {
	return MagentaString(s)
}

func SystemString(s string) string {
	return MagentaString(s)
}

func PendingString(s string) string {
	return CyanString(s)
}

func CommitString(s string) string {
	return GreenString(s)
}

func BlockString(s string) string {
	return YellowString(s)
}

func FatalString(s string) string {
	return RedString(s)
}

func TestInfoString(s string) string {
	return CyanString(s)
}

func TestOracleString(s string) string {
	return MagentaString(s)
}

func TestDecoratorString(s string) string {
	return GreenString(s)
}
