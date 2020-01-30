package world

import (
	"aakimov/marsgame/helpers"
	"aakimov/marsgame/server"
	"encoding/json"
	"log"
	"math/rand"
	"time"
)

// max moving forward per turn
const MaxMovingLength float64 = 7

// max rotation per turn in radians
const MaxRotationValue float64 = 0.1
const MaxCannonRotationValue float64 = 0.8
const MissileSpeed = 50

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
		width:          3000,
		height:         2000,
		runSpeedMs:     100,
		codeRunSpeedMs: 1000,
	}
}

const RandObjNum = 10

func prettyPrint(msg string, obj interface{}) {
	str, _ := json.MarshalIndent(obj, "", "   ")
	log.Println(msg, string(str))
}

func (w *World) MakeRandomObjects() {
	for i := 0; i < RandObjNum; i++ {
		x := float64(rand.Int31n(int32(w.width-800))) + 200.
		y := float64(rand.Int31n(int32(w.height-500))) + 200.
		//X := 1500.
		//Y := 1500.
		w.objCount += 1
		newObj := &Object{
			Id:              w.objCount,
			Type:            TypeRock,
			Pos:             Point{X: x, Y: y},
			CollisionRadius: 100,
		}
		w.objects[w.objCount] = newObj
	}
}

func (w *World) sendWorldInit(p *Player) {
	changeByTime := NewChangeByTime(0)
	for i := 1; i <= w.objCount; i++ {
		o := w.objects[i]
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

func (w *World) Run() {
	ticker := time.NewTicker(w.runSpeedMs * time.Millisecond)
	go w.sendChangelogLoop()

	serverStartTime := time.Now()
	lastTime := serverStartTime

	//log.Printf("start %v\n", serverStartTime)
	// endless loop here
	for t := range ticker.C {
		w.timeId = helpers.TimeStampDif(serverStartTime, t)
		timeDelta := time.Since(lastTime)
		lastTime = t
		//log.Printf("Game tick %v\n", t)

		w.listenChannels()
		changeByTime := NewChangeByTime(w.timeId)
		for _, player := range w.players {
			if ch := player.run(timeDelta); ch != nil {
				changeByTime.Add(ch)
			}
		}
		for id, object := range w.objects {
			if ch := object.run(w); ch != nil {
				for id1, object1 := range w.objects {
					if id1 == id {
						continue
					}
					if object.isCollideWith(object1) && object.getType() == TypeMissile {
						ch.Delete = true
						ch.DeleteOtherId = object1.getId()
						delete(w.objects, object1.getId())
						break
					}
				}
				changeByTime.Add(ch)
				if ch.Delete {
					delete(w.objects, id)
				}
			}
		}

		if changeByTime.IsNotEmpty() {
			w.changeLog.AddToBuffer(changeByTime)
		}
	}
}

func (w *World) listenChannels() {
	for {
		select {
		case client := <-w.Server.NewClientCh:
			player := NewPlayer(client.Id, client, w, w.codeRunSpeedMs)
			log.Printf("New player [%d] added to the game", player.id)

			w.players[player.id] = player
			w.sendWorldInit(player)
			go player.runProgram()
			go player.listen()
		case saveCode := <-w.Server.SaveAstCodeCh:
			player, ok := w.players[saveCode.UserId]
			if !ok {
				log.Fatalf("Save code attempt for inexistant player [%d]", saveCode.UserId)
			}
			player.saveAstCode(saveCode.SourceCode)
		case programFlowCmd := <-w.Server.ProgramFlowCh:
			player, ok := w.players[programFlowCmd.UserId]
			if !ok {
				log.Fatalf("Save code attempt for inexistant player [%d]", programFlowCmd.UserId)
			}
			player.operateState(programFlowCmd.FlowCmd)
		case o := <-w.newObjectsCh:
			w.objCount += 1
			o.setId(w.objCount)
			w.objects[w.objCount] = o
		default:
			return
		}
	}
}

func (w *World) sendChangelogLoop() {
	for {
		select {
		case ch := <-w.changeLog.changesByTimeCh:
			if w.changeLog.AddAndCheckSize(ch) {
				w.changeLog.Optimize()
				command := server.PackStructToCommand("worldChanges", w.changeLog.changesByTimeLog)
				for _, player := range w.players {
					player.client.SendCommand(command)
				}
				w.changeLog.changesByTimeLog = make([]*ChangeByTime, 0, ChangelogBufferSize)
			}
		}
	}
}
