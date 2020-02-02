package world

import (
	"aakimov/marsgame/server"
	"math/rand"
	"time"
)

// max moving forward per turn
const MaxMovingLength float64 = 7

// max rotation per turn in radians
const MaxRotationValue float64 = 0.1
const MaxCannonRotationValue float64 = 0.8
const MissileSpeed = 50
const WorldWide = 30000

type World struct {
	Server         *server.Server
	players        map[int]*Player
	objects        map[int]IObject
	changeLog      *ChangeLog
	timeId         int64
	objCount       int
	newObjectsCh   chan IObject
	width          int
	height         int
	runSpeedMs     time.Duration
	codeRunSpeedMs time.Duration
}

func NewWorld(server *server.Server) World {
	return World{
		Server:         server,
		players:        make(map[int]*Player),
		objects:        make(map[int]IObject),
		changeLog:      NewChangeLog(),
		newObjectsCh:   make(chan IObject, 10),
		width:          WorldWide,
		height:         WorldWide,
		runSpeedMs:     100,
		codeRunSpeedMs: 1000,
	}
}

const RandRocksNum = 50
const RandEnemyMechNum = 10

func (w *World) MakeRandomObjects() {
	for i := 0; i < RandRocksNum; i++ {
		x := 10000.
		y := x
		for x > 9800 && x < 10200 && y > 9800 && y < 10200 {
			x = float64(rand.Int31n(8000)) + 6000.
			y = float64(rand.Int31n(8000)) + 6000.
		}
		w.objCount += 1
		newObj := &Object{
			Id:              w.objCount,
			Type:            TypeRock,
			Pos:             Point{X: x, Y: y},
			CollisionRadius: 100,
		}
		w.objects[w.objCount] = newObj
	}
	for i := 0; i < RandEnemyMechNum; i++ {
		x := 10000.
		y := x
		for x > 9800 && x < 10200 && y > 9800 && y < 10200 {
			x = float64(rand.Int31n(8000)) + 6000.
			y = float64(rand.Int31n(8000)) + 6000.
		}
		w.objCount += 1
		newObj := &Object{
			Id:              w.objCount,
			Type:            TypeEnemyMech,
			Pos:             Point{X: x, Y: y},
			CollisionRadius: 100,
			Speed:           rand.Float64()*50 + 5.,
			AngleSpeed:      rand.Float64()*1.2 - 0.7,
		}
		w.objects[w.objCount] = newObj
	}
}

func (w *World) sendWorldInit(p *Player) {
	changeByTime := NewChangeByTime(0)
	for _, p := range w.players {
		changeByTime.Add(&ChangeByObject{
			ObjType: TypePlayer,
			ObjId:   p.id,
			Pos:     &p.mech.Pos,
			Angle:   &p.mech.Angle,
		})
	}
	for _, o := range w.objects {
		pos := o.getPos()
		changeByTime.Add(&ChangeByObject{
			ObjType: o.getType(),
			ObjId:   o.getId(),
			Pos:     &pos,
		})
	}
	ch := NewChangeLog()
	ch.AddAndCheckSize(changeByTime)

	command := server.PackStructToCommand("worldInit", ch.changesByTimeLog)
	p.client.SendCommand(command)
}
