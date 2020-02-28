package wal

import (
	"aakimov/marsgame/flatbuffers/WalBuffers"
	"aakimov/marsgame/server"
)
import flatbuffers "github.com/google/flatbuffers/go"

type Sender struct {
	terminateCh   chan bool
	logCh         chan *Log
	subscribeCh   chan *server.Client
	unsubscribeCh chan *server.Client
	clients       map[uint32]*server.Client
}

func NewSender() *Sender {
	return &Sender{
		terminateCh:   make(chan bool, 1),
		logCh:         make(chan *Log, 10),
		subscribeCh:   make(chan *server.Client, 10),
		unsubscribeCh: make(chan *server.Client, 10),
		clients:       make(map[uint32]*server.Client),
	}
}

func (s *Sender) Subscribe(client *server.Client) {
	s.subscribeCh <- client
}

func (s *Sender) Unsubscribe(client *server.Client) {
	s.unsubscribeCh <- client
}

func (s *Sender) Terminate() {
	s.terminateCh <- true
}

func (s *Sender) SendLoop() {
	for {
		select {
		case <-s.terminateCh:
			return
		case client := <-s.subscribeCh:
			s.clients[client.Id] = client
		case client := <-s.unsubscribeCh:
			delete(s.clients, client.Id)
		case log := <-s.logCh:
			buf := s.logToBuffer(log)
			for _, client := range s.clients {
				client.SendBuffer(buf)
			}
		}
	}
}

func (s *Sender) logToBuffer(logToBuff *Log) []byte {
	builder := flatbuffers.NewBuilder(1024)
	WalBuffers.LogStartTimeIdsVector(builder, len(logToBuff.TimeIds))
	for _, v := range logToBuff.TimeIds {
		builder.PrependInt32(int32(v))
	}
	timeIdsBuffObj := builder.EndVector(len(logToBuff.TimeIds))

	objsCount := len(logToBuff.Objects)
	objectLogBuffers := make([]flatbuffers.UOffsetT, objsCount)

	i := 0
	for _, obj := range logToBuff.Objects {
		timeLogsCount := len(obj.Times)
		timeLogBuffers := make([]flatbuffers.UOffsetT, timeLogsCount)
		for i, timeLog := range obj.Times {
			didCount := len(timeLog.DeleteOtherObjectIds)
			WalBuffers.TimeLogStartDeleteOtherIdsVector(builder, didCount)
			for _, did := range timeLog.DeleteOtherObjectIds {
				builder.PrependUint32(did)
			}
			didsBuffObject := builder.EndVector(didCount)

			WalBuffers.TimeLogStart(builder)
			WalBuffers.TimeLogAddTimeId(builder, int32(timeLog.TimeId))
			WalBuffers.TimeLogAddX(builder, int32(timeLog.X))
			WalBuffers.TimeLogAddY(builder, int32(timeLog.Y))
			WalBuffers.TimeLogAddAngle(builder, float32(timeLog.Angle))
			WalBuffers.TimeLogAddCannonAngle(builder, float32(timeLog.CannonAngle))
			WalBuffers.TimeLogAddCannonRotation(builder, float32(timeLog.CannonRotation))
			WalBuffers.TimeLogAddCannonUntilTimeId(builder, int32(timeLog.CannonUntilTimeId))
			WalBuffers.TimeLogAddFire(builder, timeLog.Fire)
			WalBuffers.TimeLogAddIsDelete(builder, timeLog.Delete)
			WalBuffers.TimeLogAddVelocityX(builder, float32(timeLog.VelocityX))
			WalBuffers.TimeLogAddVelocityY(builder, float32(timeLog.VelocityY))
			WalBuffers.TimeLogAddVelocityRotation(builder, float32(timeLog.VelocityRotation))
			WalBuffers.TimeLogAddVelocityUntilTimeId(builder, int32(timeLog.VelocityUntilTimeId))
			WalBuffers.TimeLogAddDeleteOtherIds(builder, didsBuffObject)
			timeLogBuffers[i] = WalBuffers.TimeLogEnd(builder)
		}
		WalBuffers.ObjectLogStartTimesVector(builder, timeLogsCount)
		for _, buffer := range timeLogBuffers {
			builder.PrependUOffsetT(buffer)
		}
		timeLogBuffersObject := builder.EndVector(timeLogsCount)

		WalBuffers.ObjectLogStart(builder)
		WalBuffers.ObjectLogAddId(builder, obj.Id)
		WalBuffers.ObjectLogAddObjectType(builder, obj.ObjType)
		WalBuffers.ObjectLogAddTimes(builder, timeLogBuffersObject)
		objectLogBuffers[i] = WalBuffers.ObjectLogEnd(builder)
		i++
	}

	WalBuffers.LogStartObjectsVector(builder, objsCount)
	for _, buffer := range objectLogBuffers {
		builder.PrependUOffsetT(buffer)
	}
	objectLogBuffersObject := builder.EndVector(objsCount)

	WalBuffers.LogStart(builder)
	WalBuffers.LogAddTimeIds(builder, timeIdsBuffObj)
	WalBuffers.LogAddObjects(builder, objectLogBuffersObject)
	logBufferObj := WalBuffers.LogEnd(builder)

	builder.Finish(logBufferObj)
	buf := builder.FinishedBytes()
	return buf
}
