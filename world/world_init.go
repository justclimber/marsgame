package world

import (
	"aakimov/marsgame/flatbuffers/CommandsBuffer"
	"aakimov/marsgame/flatbuffers/InitBuffers"
	"aakimov/marsgame/physics"
	"aakimov/marsgame/server"
	"aakimov/marsgame/wal"
	flatbuffers "github.com/google/flatbuffers/go"
	"math/rand"
)

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
			x = float64(rand.Int31n(8000)) + 6000.
			y = float64(rand.Int31n(8000)) + 6000.
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
			obj.Speed = rand.Float64()*300 + 50.
			obj.AngleSpeed = rand.Float64()*1.2 - 0.7
		}},
	} {
		w.MakeRandomObjectsByType(v)
	}
}

func (w *World) createPlayerAndBootstrap(client *server.Client) *Player {
	if w.timer.IsStopped() {
		w.timer.Start(w.stopByTimer)
	}
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
	w.sendInitDataToNewPlayer(player)
	go player.programLoop()
	go player.listen()

	return player
}

func (w *World) sendInitDataToNewPlayer(player *Player) {
	builder := flatbuffers.NewBuilder(1024)
	InitBuffers.TimerStart(builder)
	InitBuffers.TimerAddState(builder, w.timer.State())
	InitBuffers.TimerAddValue(builder, int32(w.timer.Value().Seconds()))
	timerBufferObj := InitBuffers.TimerEnd(builder)

	InitBuffers.InitStart(builder)
	InitBuffers.InitAddTimer(builder, timerBufferObj)
	initBufferObj := InitBuffers.InitEnd(builder)
	builder.Finish(initBufferObj)
	buf := builder.FinishedBytes()

	buf = append([]byte{byte(CommandsBuffer.CommandInit)}, buf...)
	player.client.SendBuffer(buf)
}

func (w *World) stopByTimer() {
	// @todo: send commands to clients with stopping the mage
	w.stop()
}

func (w *World) stop() {
	w.wal.Sender.Terminate()
	for _, p := range w.players {
		p.flowCh <- Terminate
	}
}

func (w *World) reset() {
	w.stop()
	w.wal = wal.NewWal()
	for i, p := range w.players {
		astProgram := p.mainProgram.astProgram
		w.players[i] = w.createPlayerAndBootstrap(p.client)
		w.players[i].mainProgram.astProgram = astProgram
	}
	w.objects = make(map[uint32]IObject)
	w.objCount = 0
	w.MakeRandomObjects()
	go w.wal.Sender.SendLoop()
}
