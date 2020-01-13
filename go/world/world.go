package world

import (
	"aakimov/marsgame/go/server"
	"log"
	"time"
)

// max moving forward per turn
const MaxMovingLength float64 = 7

// max rotation per turn in radians
const MaxRotationValue float64 = 0.5

type World struct {
	Server    *server.Server
	players   map[string]*Player
	objects   map[string]*Object
	changeLog *ChangeLog
}

func NewWorld(server *server.Server) World {
	return World{
		Server:    server,
		players:   make(map[string]*Player),
		objects:   make(map[string]*Object),
		changeLog: NewChangeLog(),
	}
}

func makeTimestamp(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func timeStampDif(t1, t2 time.Time) int64 {
	return makeTimestamp(t2) - makeTimestamp(t1)
}

func (w *World) codeRun() {
	// players are empty at start, so this block is for future
	for _, player := range w.players {
		go player.mainProgram.Run()
		go player.listen()
	}
}

func (w *World) Run() {
	ticker := time.NewTicker(200 * time.Millisecond)
	go w.sendChangelogLoop()
	w.codeRun()

	serverStartTime := time.Now()

	// endless loop here
	for t := range ticker.C {
		timeId := timeStampDif(serverStartTime, t)
		log.Printf("Game tick %v\n", timeId)

		w.listenChannels()
		changeByTime := NewChangeByTime(timeId)
		for _, player := range w.players {
			if ch := w.runPlayer(player); ch != nil {
				changeByTime.Add(ch)
			}
		}
		for _, object := range w.objects {
			w.runObject(object)
		}
		if changeByTime.IsNotEmpty() {
			w.changeLog.AddToBuffer(changeByTime)
		}
	}
}

func (w *World) listenChannels() {
	select {
	case client := <-w.Server.NewClientCh:
		player := NewPlayer(client.Id, client, NewMech())
		log.Printf("New player [%s] added to the game", player.id)

		w.players[player.id] = player
		go player.mainProgram.Run()
		go player.listen()
	case saveCode := <-w.Server.SaveAstCodeCh:
		player, ok := w.players[saveCode.UserId]
		if !ok {
			log.Fatalf("Save code attempt for inexistant player [%s]", saveCode.UserId)
		}
		player.saveAstCode(saveCode.SourceCode)
	default:
		// noop
	}
}

func (w *World) runPlayer(player *Player) *ChangeByObject {
	mech := &player.mech
	changeByObject := ChangeByObject{
		ObjType: TypePlayer,
		ObjId:   player.id,
	}
	mech.mu.Lock()
	if mech.RotateThrottle != 0 {
		mech.Object.Angle += mech.RotateThrottle * MaxRotationValue
		changeByObject.Angle = mech.Object.Angle
	}
	if mech.Throttle != 0 {
		length := mech.Throttle * MaxMovingLength
		mech.Object.Pos.MoveForward(mech.Object.Angle, length)
		changeByObject.Pos = mech.Object.Pos
	}
	mech.mu.Unlock()

	if mech.RotateThrottle != 0 || mech.Throttle != 0 {
		return &changeByObject
	}
	return nil
}

func (w *World) runObject(object *Object) {

}

func (w *World) sendChangelogLoop() {
	for {
		select {
		case ch := <-w.changeLog.changesByTimeCh:
			if w.changeLog.AddAndCheckSize(ch) {
				w.changeLog.Optimize()
				command := PackChangesToCommand(w.changeLog.changesByTimeLog)
				for _, player := range w.players {
					player.client.SendCommand(command)
				}
				w.changeLog.changesByTimeLog = make([]*ChangeByTime, 0, ChangelogBufferSize)
			}
		default:
			// noop
		}
	}
}
