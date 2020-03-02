package world

import (
	"aakimov/marsgame/server"
	"aakimov/marsgame/timer"
	"aakimov/marsgame/wal"

	"time"
)

const Wide = 30000

type World struct {
	Server         *server.Server
	players        map[uint32]*Player
	objects        map[uint32]IObject
	objCount       uint32
	newObjectsCh   chan IObject
	width          int
	height         int
	runSpeedMs     time.Duration
	codeRunSpeedMs time.Duration
	wal            *wal.Wal
	timer          *timer.Timer
}

func NewWorld(server *server.Server) World {
	return World{
		Server:         server,
		players:        make(map[uint32]*Player),
		objects:        make(map[uint32]IObject),
		newObjectsCh:   make(chan IObject, 10),
		width:          Wide,
		height:         Wide,
		runSpeedMs:     100,
		codeRunSpeedMs: 1000,
		wal:            wal.NewWal(),
	}
}

func ObjectTypeToInt(objType string) int8 {
	var objTypeToIntMap = map[string]int8{
		TypePlayer:    0,
		TypeEnemyMech: 1,
		TypeRock:      2,
		TypeXelon:     3,
		TypeMissile:   4,
	}
	return objTypeToIntMap[objType]
}
