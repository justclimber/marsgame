package world

import (
	"aakimov/marsgame/helpers"
	"aakimov/marsgame/server"
	"log"
	"math/rand"
	"time"
)

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
			x := float64(rand.Int31n(1000)) + 9500.
			y := float64(rand.Int31n(1000)) + 9500.
			mech := NewMech(x, y)
			player := NewPlayer(client.Id, client, w, mech, w.codeRunSpeedMs)
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
