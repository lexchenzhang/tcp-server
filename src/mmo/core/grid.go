package core

import (
	"fmt"
	"sync"
)

type Grid struct {
	GID       int
	MinX      int
	MaxX      int
	MinY      int
	MaxY      int
	playerIDs map[int]bool
	pIDLock   sync.RWMutex
}

func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

func (g *Grid) Add(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()
	g.playerIDs[playerID] = true
}

func (g *Grid) Remove(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()
	delete(g.playerIDs, playerID)
}

func (g *Grid) GetPlayerIDs() (playerIDs []int) {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()
	for k := range g.playerIDs {
		playerIDs = append(playerIDs, k)
	}
	return playerIDs
}

func (g *Grid) String() string {
	return fmt.Sprintf("Grid id: %d, minX:%d, maxX:%d, minY:%d, maxY:%d, playerIDs:%v", g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
