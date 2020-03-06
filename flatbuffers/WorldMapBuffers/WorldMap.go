// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package WorldMapBuffers

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type WorldMap struct {
	_tab flatbuffers.Table
}

func GetRootAsWorldMap(buf []byte, offset flatbuffers.UOffsetT) *WorldMap {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &WorldMap{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *WorldMap) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *WorldMap) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *WorldMap) Layers(obj *TileLayer, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *WorldMap) LayersLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *WorldMap) Width() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *WorldMap) MutateWidth(n int32) bool {
	return rcv._tab.MutateInt32Slot(6, n)
}

func (rcv *WorldMap) Height() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *WorldMap) MutateHeight(n int32) bool {
	return rcv._tab.MutateInt32Slot(8, n)
}

func WorldMapStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func WorldMapAddLayers(builder *flatbuffers.Builder, layers flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(layers), 0)
}
func WorldMapStartLayersVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func WorldMapAddWidth(builder *flatbuffers.Builder, width int32) {
	builder.PrependInt32Slot(1, width, 0)
}
func WorldMapAddHeight(builder *flatbuffers.Builder, height int32) {
	builder.PrependInt32Slot(2, height, 0)
}
func WorldMapEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}