package world

import (
	"aakimov/marsgame/physics"
	"aakimov/marsgame/server"
	"aakimov/marsgame/wal"
	"math/rand"
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

func (w *World) Bootstrap() {
	w.MakeRandomObjects()
	go w.run()
	go w.wal.Sender.SendLoop()
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
			x = float64(rand.Int31n(10000)) + 6000.
			y = float64(rand.Int31n(10000)) + 6000.
		}
		w.objCount += 1
		newObj := &Object{
			Obj: physics.Obj{
				Id:              w.objCount,
				Type:            seed.objType,
				Pos:             physics.Point{X: x, Y: y},
				CollisionRadius: seed.collisionRadius,
				Velocity:        &physics.Vector{},
				Direction:       physics.MakeNormalVectorByAngle(0),
			},
			wal: w.wal.NewObjectObserver(w.objCount, ObjectTypeToInt(seed.objType)),
		}

		if seed.extraCallback != nil {
			seed.extraCallback(newObj)
		}
		newObj.wal.AddAngle(newObj.Angle)
		newObj.wal.AddPosAndVelocityLen(newObj.Pos, newObj.Speed)
		newObj.wal.AddRotation(0)
		w.objects[w.objCount] = newObj
	}
}

func (w *World) MakeRandomObjects() {
	for _, v := range []RandomObjSeed{
		{TypeRock, 30, 100, nil},
		{TypeXelon, 30, 50, nil},
		{TypeEnemyMech, 10, 100, func(obj *Object) {
			obj.Speed = rand.Float64()*50 + 5.
			obj.AngleSpeed = rand.Float64()*1.2 - 0.7
		}},
	} {
		w.MakeRandomObjectsByType(v)
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

func (w *World) createPlayerAndBootstrap(client *server.Client) *Player {
	x := float64(rand.Int31n(1000)) + 9500.
	y := float64(rand.Int31n(1000)) + 9500.
	mech := NewMech(x, y)
	player := NewPlayer(
		client.Id,
		client,
		w,
		mech,
		w.codeRunSpeedMs,
		w.wal.NewObjectObserver(client.Id, ObjectTypeToInt(TypePlayer)),
	)
	player.wal.AddPosAndVelocityLen(mech.Pos, mech.Velocity.Len())
	player.wal.AddAngle(mech.Angle)
	player.wal.AddRotation(0)
	player.wal.AddCannonAngle(0)
	player.wal.AddRotation(0)

	w.players[player.id] = player
	w.wal.Sender.Subscribe(player.client)
	go player.runProgram()
	go player.listen()

	return player
}

func (w *World) reset() {
	w.wal.Sender.Terminate()
	w.wal = wal.NewWal()
	for i, p := range w.players {
		p.flowCh <- Terminate
		astProgram := p.mainProgram.astProgram
		w.players[i] = w.createPlayerAndBootstrap(p.client)
		w.players[i].mainProgram.astProgram = astProgram
	}
	w.objects = make(map[uint32]IObject)
	w.objCount = 0
	w.MakeRandomObjects()
	go w.wal.Sender.SendLoop()
}
