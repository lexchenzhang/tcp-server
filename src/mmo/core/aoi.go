package core

import "fmt"

type AOIManager struct {
	MinX  int
	MaxX  int
	CntsX int
	MinY  int
	MaxY  int
	CntsY int
	grids map[int]*Grid
}

func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		grids: make(map[int]*Grid),
	}

	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			gid := y*cntsX + x
			aoiMgr.grids[gid] = NewGrid(
				gid,
				aoiMgr.MinX+x*aoiMgr.gridWidth(),
				aoiMgr.MinX+(x+1)*aoiMgr.gridWidth(),
				aoiMgr.MinY+y*aoiMgr.gridHeight(),
				aoiMgr.MinY+(y+1)*aoiMgr.gridHeight())
		}
	}

	return aoiMgr
}

func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

func (m *AOIManager) gridHeight() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManagr:\nminX:%d, maxX:%d, cntsX:%d, minY:%d, maxY:%d, cntsY:%d\n GrIDs in AOI Manager:\n",
		m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)
	for _, grID := range m.grids {
		s += fmt.Sprintln(grID)
	}
	return s
}
