package world

import (
	"aakimov/marsgame/flatbuffers/generated/WalBuffers"
	"aakimov/marsgame/physics"
	"aakimov/marsgame/server"
	"aakimov/marsgame/wal"

	"log"
	"math/rand"
)

func (w *World) Bootstrap() {
	w.worldmap.Parse("worldmap/firstmap.tmx")
	w.makeObjectsFromWorldMap()
	//w.MakeRandomObjects()
	go w.run()
	go w.wal.Sender.SendLoop()
}

type RandomObjSeed struct {
	objType         string
	count           int
	collisionRadius int
	extraCallback   func(*Object)
}

func (w *World) makeObjectsFromWorldMap() {
	var newObj *Object
	for _, entity := range w.worldmap.Entities {
		switch entity.EntityType {
		case WalBuffers.ObjectTypexelon, WalBuffers.ObjectTypespore:
			newObj = &Object{
				Obj: physics.Obj{
					Id:              entity.Id,
					Type:            EntityTypeToObjectType(entity.EntityType),
					Pos:             entity.Pos,
					CollisionRadius: 20,
					Velocity:        &physics.Vector{},
					Direction:       physics.MakeNormalVectorByAngle(0),
				},
				wal: w.wal.NewObjectObserver(entity.Id, ObjectTypeToInt(EntityTypeToObjectType(entity.EntityType))),
			}
			w.objects[entity.Id] = newObj
			w.objCount++
			newObj.wal.AddAngle(newObj.Angle)
			newObj.wal.AddPosAndVelocityLen(newObj.Pos, newObj.Speed)
			newObj.wal.AddRotation(0)
		case WalBuffers.ObjectTypeplayer:
			w.playerPosSlots = append(w.playerPosSlots, entity.Pos)
		default:
			log.Fatalf("Unsopported object type %v", entity.EntityType)
		}
	}
}

func (w *World) MakeRandomObjectsByType(seed RandomObjSeed) {
	for i := 0; i < seed.count; i++ {
		x := 10000.
		y := x
		for x > 9800 && x < 10200 && y > 9800 && y < 10200 {
			x = float64(rand.Int31n(2000)) + 9000.
			y = float64(rand.Int31n(2000)) + 9000.
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
		{TypeRock, 30, 30, nil},
		{TypeXelon, 30, 10, nil},
		{TypeEnemyMech, 10, 20, func(obj *Object) {
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
	if len(w.playerPosSlots) == 0 {
		log.Fatalln("No free slots for new player, exiting")
	}
	playerPosSlot := w.playerPosSlots[0]
	mech := NewMech(playerPosSlot.X, playerPosSlot.Y)
	w.playerPosSlots = w.playerPosSlots[1:]
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
	player.client.SendBuffer(w.initDataToBuf())
}

func (w *World) stopByTimer() {
	// @todo: send commands to clients with stopping the game
	w.stop()
	log.Fatalln("World stopped by timer")
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
