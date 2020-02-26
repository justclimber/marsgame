package world

import (
	"log"
	"time"
)

func (w *World) run() {
	ticker := time.NewTicker(w.runSpeedMs * time.Millisecond)
	serverStartTime := time.Now()
	lastTime := serverStartTime

	//log.Printf("start %v\n", serverStartTime)

	// world mechanics processing loop
	for t := range ticker.C {
		timeId := t.Sub(serverStartTime).Milliseconds()
		timeDelta := t.Sub(lastTime)
		lastTime = t
		//log.Printf("Game tick %v\n", t)
		//log.Printf("Time delta %v\n", timeDelta.Milliseconds())

		w.listenChannels()
		for _, player := range w.players {
			player.run(timeDelta, timeId)
		}
		for _, object := range w.objects {
			object.run(w, timeDelta, timeId)
		}
		w.wal.Commit(timeId)
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
			o.setObjectManager(w.wal.NewObjectObserver(w.objCount, o.getType()))
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
