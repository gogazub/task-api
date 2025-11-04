package model

// Result - result of execution of code. Contains text from stderr and stdout
type Result struct {
	Error  []byte `db:"stderr"`
	Output []byte `db:"stdout"`
}
