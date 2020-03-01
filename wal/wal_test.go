package wal

import (
	"aakimov/marsgame/helpers"
	"aakimov/marsgame/physics"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWalSimple(t *testing.T) {
	mainWal := NewWal()
	objectLog := mainWal.NewObjectObserver(1, 1)
	objectLog.AddPosAndVelocityLen(physics.Point{X: 10, Y: 10}, 0)
	objectLog.Commit(101)
	mainWal.Commit(101)
	log := mainWal.logBuffer

	require.Equal(t, `{
   "TimeIds": [
      101
   ],
   "Objects": [
      {
         "Id": 1,
         "ObjType": 1,
         "Times": [
            {
               "TimeId": 101,
               "IsNew": false,
               "X": 10,
               "Y": 10,
               "Angle": 99999999,
               "CannonAngle": 99999999,
               "CannonRotation": 99999999,
               "CannonUntilTimeId": 99999999,
               "Fire": false,
               "Delete": false,
               "Explode": false,
               "ExplodeOther": false,
               "DeleteOtherObjectIds": null,
               "VelocityX": 0,
               "VelocityY": 0,
               "VelocityRotation": 99999999,
               "VelocityUntilTimeId": 99999999
            }
         ]
      }
   ]
}`, helpers.Pretty(log))
}

func TestWal2SameSimpleRecords(t *testing.T) {
	mainWal := NewWal()
	objectLog := mainWal.NewObjectObserver(1, 1)
	objectLog.AddPosAndVelocityLen(physics.Point{X: 10, Y: 10}, 0)
	objectLog.Commit(101)
	mainWal.Commit(101)
	objectLog.AddPosAndVelocityLen(physics.Point{X: 10, Y: 10}, 0)
	objectLog.Commit(202)
	mainWal.Commit(202)
	log := mainWal.logBuffer

	require.Equal(t, `{
   "TimeIds": [
      101,
      202
   ],
   "Objects": [
      {
         "Id": 1,
         "ObjType": 1,
         "Times": [
            {
               "TimeId": 101,
               "IsNew": false,
               "X": 10,
               "Y": 10,
               "Angle": 99999999,
               "CannonAngle": 99999999,
               "CannonRotation": 99999999,
               "CannonUntilTimeId": 99999999,
               "Fire": false,
               "Delete": false,
               "Explode": false,
               "ExplodeOther": false,
               "DeleteOtherObjectIds": null,
               "VelocityX": 0,
               "VelocityY": 0,
               "VelocityRotation": 99999999,
               "VelocityUntilTimeId": 202
            }
         ]
      }
   ]
}`, helpers.Pretty(log))
}

//func TestWal2AlmostSameRecords(t *testing.T) {
//	mainWal := NewWal()
//	objectLog := mainWal.NewObjectObserver(1, 1)
//	objectLog.AddPosAndVelocityLen(physics.Point{X: 10, Y: 10}, &physics.Vector{})
//	objectLog.Commit(101)
//	mainWal.Commit(101)
//	objectLog.AddPosAndVelocityLen(physics.Point{X: 10, Y: 10}, &physics.Vector{})
//	objectLog.AddCannonAngle(1)
//	objectLog.Commit(202)
//	mainWal.Commit(202)
//	log := mainWal.logBuffer
//
//	require.Equal(t, `{
//   "TimeIds": [
//      101,
//      202
//   ],
//   "Objects": [
//      {
//         "Id": 1,
//         1: 1,
//         "l": [
//            {
//               "Id": 101,
//               "x": 10,
//               "y": 10,
//               "vx": 0,
//               "vy": 0,
//               "vl": 0,
//               "vu": 202
//            },
//            {
//               "Id": 202,
//               "ca": 1,
//            }
//         ]
//      }
//   ]
//}`, helpers.Pretty(log))
//}
