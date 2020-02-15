package world

import (
	"aakimov/marsgame/helpers"
	"aakimov/marsgame/server"
	"log"
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
			player := w.createPlayerAndBootstrap(client)
			w.sendWorldInit(player)
			log.Printf("New player [%d] added to the game", player.id)
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
		case c := <-w.Server.CommandsCh:
			player, ok := w.players[c.UserId]
			if !ok {
				log.Fatalf("Command from inexistant player [%d]", c.UserId)
			}
			switch c.Command.Type {
			case "resetMech":
				player.flowCh <- Terminate
				player := w.createPlayerAndBootstrap(player.client)
				log.Printf("Player [%d] has ben reset\n", player.id)
			case "resetWorld":
				w.reset()
				log.Println("World has ben reset")
			default:
				log.Fatalf("Command '%s' not supported", c.Command.Type)
			}
		default:
			return
		}
	}
}

func (w *World) sendChangelogLoop() {
	for {
		select {
		case <-w.changeLog.terminateCh:
			return
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
