package world

import (
	"github.com/justclimber/marsgame/flatbuffers/generated/WalBuffers"
	"github.com/justclimber/marsgame/physics"
	"github.com/justclimber/marsgame/server"
	"github.com/justclimber/marsgame/timer"
	"github.com/justclimber/marsgame/wal"
	"github.com/justclimber/marsgame/worldmap"

	"time"
)

const Wide = 30000
const TimeMultiplicator = 1

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
	worldmap       *worldmap.WorldMap
	playerPosSlots []physics.Point
	lastNewObjId   uint32
}

func NewWorld(server *server.Server) World {
	return World{
		Server:         server,
		players:        make(map[uint32]*Player),
		objects:        make(map[uint32]IObject),
		newObjectsCh:   make(chan IObject, 10),
		width:          Wide,
		height:         Wide,
		runSpeedMs:     100 / TimeMultiplicator,
		codeRunSpeedMs: 1000 / TimeMultiplicator,
		wal:            wal.NewWal(),
		timer:          timer.NewTimer(time.Second*50, TimeMultiplicator),
		worldmap:       worldmap.NewWorldMap(),
		playerPosSlots: make([]physics.Point, 0),
		lastNewObjId:   1000000,
	}
}

func ObjectTypeToInt(objType string) int8 {
	var objTypeToIntMap = map[string]int8{
		TypePlayer:    0,
		TypeEnemyMech: 1,
		TypeRock:      2,
		TypeXelon:     3,
		TypeMissile:   4,
		TypeSpore:     5,
	}
	return objTypeToIntMap[objType]
}

func EntityTypeToObjectType(objType WalBuffers.ObjectType) string {
	var entityTypeMap = map[WalBuffers.ObjectType]string{
		WalBuffers.ObjectTypexelon: TypeXelon,
		WalBuffers.ObjectTypespore: TypeSpore,
	}
	return entityTypeMap[objType]
}

func getObjectsMeta() map[string]ObjectMeta {
	return map[string]ObjectMeta{
		TypePlayer:    {ObjectTypeToInt(TypePlayer), 25},
		TypeEnemyMech: {ObjectTypeToInt(TypeEnemyMech), 10},
		TypeRock:      {ObjectTypeToInt(TypeRock), 10},
		TypeXelon:     {ObjectTypeToInt(TypeXelon), 14},
		TypeMissile:   {ObjectTypeToInt(TypeMissile), 6},
		TypeSpore:     {ObjectTypeToInt(TypeSpore), 40},
	}
}
