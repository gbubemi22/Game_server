package game

import (
	"net"
	"sync"
)

type Player struct {
	ID     string
	Name   string
	X, Y   float64
	Health int
	Conn   net.Conn
}

type GameRoom struct {
	ID      string
	Players map[string]*Player
	Mu      sync.Mutex
	State   string // "waiting", "active", "finished"
}
