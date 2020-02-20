package wal

import (
	"aakimov/marsgame/helpers"
	"aakimov/marsgame/physics"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWallSimple(t *testing.T) {
	mainWal := NewManager()
	objectLog := mainWal.CreateObjectManager(1, "t")
	objectLog.AddPosAndVelocity(physics.Point{X: 10, Y: 10}, &physics.Vector{})
	objectLog.Commit(101)
	mainWal.Commit(101)
	log := mainWal.MainLog

	require.Equal(t, `{
   "tIds": [
      101
   ],
   "objs": [
      {
         "i": 1,
         "t": "t",
         "l": [
            {
               "i": 101,
               "x": 10,
               "y": 10,
               "vx": 0,
               "vy": 0,
               "vl": 0
            }
         ]
      }
   ]
}`, helpers.Pretty(log))
}

func TestWall2SameSimpleRecords(t *testing.T) {
	mainWal := NewManager()
	objectLog := mainWal.CreateObjectManager(1, "t")
	objectLog.AddPosAndVelocity(physics.Point{X: 10, Y: 10}, &physics.Vector{})
	objectLog.Commit(101)
	mainWal.Commit(101)
	objectLog.AddPosAndVelocity(physics.Point{X: 10, Y: 10}, &physics.Vector{})
	objectLog.Commit(202)
	mainWal.Commit(202)
	log := mainWal.MainLog

	require.Equal(t, `{
   "tIds": [
      101,
      202
   ],
   "objs": [
      {
         "i": 1,
         "t": "t",
         "l": [
            {
               "i": 101,
               "x": 10,
               "y": 10,
               "vx": 0,
               "vy": 0,
               "vl": 0,
               "vu": 202
            }
         ]
      }
   ]
}`, helpers.Pretty(log))
}

//func TestWall2AlmostSameRecords(t *testing.T) {
//	mainWal := NewManager()
//	objectLog := mainWal.CreateObjectManager(1, "t")
//	objectLog.AddPosAndVelocity(physics.Point{X: 10, Y: 10}, &physics.Vector{})
//	objectLog.Commit(101)
//	mainWal.Commit(101)
//	objectLog.AddPosAndVelocity(physics.Point{X: 10, Y: 10}, &physics.Vector{})
//	objectLog.AddCannonAngle(1)
//	objectLog.Commit(202)
//	mainWal.Commit(202)
//	log := mainWal.MainLog
//
//	require.Equal(t, `{
//   "tIds": [
//      101,
//      202
//   ],
//   "objs": [
//      {
//         "i": 1,
//         "t": "t",
//         "l": [
//            {
//               "i": 101,
//               "x": 10,
//               "y": 10,
//               "vx": 0,
//               "vy": 0,
//               "vl": 0,
//               "vu": 202
//            },
//            {
//               "i": 202,
//               "ca": 1,
//            }
//         ]
//      }
//   ]
//}`, helpers.Pretty(log))
//}
