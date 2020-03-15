// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package InitBuffers

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type ObjectMeta struct {
	_tab flatbuffers.Table
}

func GetRootAsObjectMeta(buf []byte, offset flatbuffers.UOffsetT) *ObjectMeta {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &ObjectMeta{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *ObjectMeta) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ObjectMeta) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *ObjectMeta) ObjectType() int8 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt8(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *ObjectMeta) MutateObjectType(n int8) bool {
	return rcv._tab.MutateInt8Slot(4, n)
}

func (rcv *ObjectMeta) CollisionRadius() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *ObjectMeta) MutateCollisionRadius(n int16) bool {
	return rcv._tab.MutateInt16Slot(6, n)
}

func ObjectMetaStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func ObjectMetaAddObjectType(builder *flatbuffers.Builder, objectType int8) {
	builder.PrependInt8Slot(0, objectType, 0)
}
func ObjectMetaAddCollisionRadius(builder *flatbuffers.Builder, collisionRadius int16) {
	builder.PrependInt16Slot(1, collisionRadius, 0)
}
func ObjectMetaEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
