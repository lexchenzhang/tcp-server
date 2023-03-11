package core

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewAOIManager(t *testing.T) {
	aoiMgr := NewAOIManager(0, 250, 5, 0, 250, 5)
	fmt.Println(aoiMgr)
}

func TestAOIManagerSuroundGrIDsByGID(t *testing.T) {
	aoiMgr := NewAOIManager(0, 250, 5, 0, 250, 5)
	var tests = []struct {
		name     string
		gridID   int
		expected []int
	}{
		{"case 1", 0, []int{1, 5, 6}},
		{"case 2", 4, []int{3, 8, 9}},
		{"case 3", 11, []int{5, 6, 7, 10, 12, 15, 16, 17}},
		{"case 4", 22, []int{16, 17, 18, 21, 23}},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			if _, ret := aoiMgr.GetSurroundGrIDsByGID(c.gridID); !reflect.DeepEqual(ret, c.expected) {
				t.Errorf("should be %v but got %v", c.expected, ret)
			}
		})
	}
}
