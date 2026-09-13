package server

type StateMachine interface {
	Apply(cmd []byte) []byte
}
