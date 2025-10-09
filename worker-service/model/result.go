package model

type Result struct {
	Error  []byte `db:"stderr"`
	Output []byte `db:"stdout"`
}
