package game

import "sync"

var (
	gameRooms = make(map[string]*GameRoom)
	roomsMu   sync.Mutex
)
func GetOrCreateRoom(roomID string) *GameRoom {
	roomsMu.Lock()
	defer roomsMu.Unlock()

	room, exists := gameRooms[roomID]
	if !exists {
		room = &GameRoom{
			ID:      roomID,
			Players: make(map[string]*Player),
			State:   "waiting",
		}
		gameRooms[roomID] = room
	}
	return room
}





func AddPlayerToRoom(roomID string, p *Player) {
	room := GetOrCreateRoom(roomID)
	room.Mu.Lock()
	defer room.Mu.Unlock()

	room.Players[p.ID] = p
}




func MovePlayer(roomID, playerID string, x, y float64) bool {
	room := GetOrCreateRoom(roomID)
	room.Mu.Lock()
	defer room.Mu.Unlock()

	player, exists := room.Players[playerID]
	if !exists {
		return false
	}
	player.X = x
	player.Y = y
	return true
}



func ApplyDamage(roomID, targetID string, amount int) bool {
	room := GetOrCreateRoom(roomID)
	room.Mu.Lock()
	defer room.Mu.Unlock()

	target, exists := room.Players[targetID]
	if !exists {
		return false
	}
	target.Health -= amount
	if target.Health <= 0 {
		target.Health = 0
		// TODO: Mark as dead, broadcast kill
	}
	return true
}
