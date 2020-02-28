// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package WalBuffers

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Log struct {
	_tab flatbuffers.Table
}

func GetRootAsLog(buf []byte, offset flatbuffers.UOffsetT) *Log {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Log{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *Log) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Log) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Log) TimeIds(j int) int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetInt32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *Log) TimeIdsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *Log) MutateTimeIds(j int, n int32) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateInt32(a+flatbuffers.UOffsetT(j*4), n)
	}
	return false
}

func (rcv *Log) Objects(obj *ObjectLog, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *Log) ObjectsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func LogStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func LogAddTimeIds(builder *flatbuffers.Builder, timeIds flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(timeIds), 0)
}
func LogStartTimeIdsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func LogAddObjects(builder *flatbuffers.Builder, objects flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(objects), 0)
}
func LogStartObjectsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func LogEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
