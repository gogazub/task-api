package model

// Result - result of execution of code. Contains text from stderr and stdout
type Result struct {
	Id string `db:"id"`
	Error  []byte `db:"stderr"`
	Output []byte `db:"stdout"`
}
