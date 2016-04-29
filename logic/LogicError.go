// LogicError
package logic

var (
	OverflowError = &LogicError{"Overflow"}
)

type LogicError struct {
	es string
}

func (le *LogicError) Error() string {
	return le.es
}
