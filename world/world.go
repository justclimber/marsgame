package world

import (
	"aakimov/marsgame/changelog"
	"aakimov/marsgame/physics"
	"aakimov/marsgame/server"
	"math/rand"
	"time"
)

// max moving forward per turn
const MaxMovingLength float64 = 7

// max rotation per turn in radians
const MaxRotationValue float64 = 0.1
const MaxCannonRotationValue float64 = 0.11
const MissileSpeed = 50
const WorldWide = 30000

type World struct {
	Server         *server.Server
	players        map[int]*Player
	objects        map[int]IObject
	changeLog      *changelog.ChangeLog
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
		changeLog:      changelog.NewChangeLog(),
		newObjectsCh:   make(chan IObject, 10),
		width:          WorldWide,
		height:         WorldWide,
		runSpeedMs:     100,
		codeRunSpeedMs: 1000,
	}
}

type RandomObjSeed struct {
	objType         string
	count           int
	collisionRadius int
	extraCallback   func(*Object)
}

func (w *World) MakeRandomObjectsByType(seed RandomObjSeed) {
	for i := 0; i < seed.count; i++ {
		x := 10000.
		y := x
		for x > 9800 && x < 10200 && y > 9800 && y < 10200 {
			x = float64(rand.Int31n(8000)) + 6000.
			y = float64(rand.Int31n(8000)) + 6000.
		}
		w.objCount += 1
		newObj := &Object{physics.Obj{
			Id:              w.objCount,
			Type:            seed.objType,
			Pos:             physics.Point{X: x, Y: y},
			CollisionRadius: seed.collisionRadius,
		}}
		if seed.extraCallback != nil {
			seed.extraCallback(newObj)
		}
		w.objects[w.objCount] = newObj
	}
}

func (w *World) MakeRandomObjects() {
	w.MakeRandomObjectsByType(RandomObjSeed{TypeRock, 50, 100, nil})
	w.MakeRandomObjectsByType(RandomObjSeed{TypeXelon, 50, 50, nil})
	w.MakeRandomObjectsByType(RandomObjSeed{TypeEnemyMech, 10, 100, func(obj *Object) {
		obj.Speed = rand.Float64()*50 + 5.
		obj.AngleSpeed = rand.Float64()*1.2 - 0.7
	}})
}

func (w *World) createPlayerAndBootstrap(client *server.Client) *Player {
	x := float64(rand.Int31n(1000)) + 9500.
	y := float64(rand.Int31n(1000)) + 9500.
	mech := NewMech(x, y)
	player := NewPlayer(client.Id, client, w, mech, w.codeRunSpeedMs)
	w.players[player.id] = player
	go player.runProgram()
	go player.listen()

	return player
}

func (w *World) reset() {
	for i, p := range w.players {
		p.flowCh <- Terminate
		astProgram := p.mainProgram.astProgram
		w.players[i] = w.createPlayerAndBootstrap(p.client)
		w.players[i].mainProgram.astProgram = astProgram
	}
	w.objects = make(map[int]IObject)
	w.objCount = 0
	w.MakeRandomObjects()
	w.changeLog.Terminate()
	w.changeLog = changelog.NewChangeLog()
	go w.sendChangelogLoop()
	for _, p := range w.players {
		w.sendWorldInit(p)
	}
}

func (w *World) sendWorldInit(p *Player) {
	changeByTime := changelog.NewChangeByTime(0)
	for _, p := range w.players {
		changeByTime.Add(&changelog.ChangeByObject{
			ObjType: TypePlayer,
			ObjId:   p.id,
			Pos:     &p.mech.Pos,
			Angle:   &p.mech.Angle,
		})
	}
	for _, o := range w.objects {
		pos := o.getPos()
		changeByTime.Add(&changelog.ChangeByObject{
			ObjType: o.getType(),
			ObjId:   o.getId(),
			Pos:     &pos,
		})
	}
	ch := changelog.NewChangeLog()
	ch.AddAndCheckSize(changeByTime)

	command := server.PackStructToCommand("worldInit", ch.GetLog())
	p.client.SendCommand(command)
}
