package world

import (
	"github.com/justclimber/marsgame/flatbuffers/generated/CommandsBuffer"
	"github.com/justclimber/marsgame/flatbuffers/generated/InitBuffers"
	"github.com/justclimber/marsgame/flatbuffers/generated/WorldMapBuffers"

	flatbuffers "github.com/google/flatbuffers/go"
)

func (w *World) initDataToBuf() []byte {
	builder := flatbuffers.NewBuilder(1024)
	InitBuffers.TimerStart(builder)
	InitBuffers.TimerAddState(builder, w.timer.State())
	InitBuffers.TimerAddValue(builder, int32(w.timer.Value().Seconds()))
	timerBufferObj := InitBuffers.TimerEnd(builder)

	objectsMeta := getObjectsMeta()
	objectsMetaCount := len(objectsMeta)
	objectsMetaBuffers := make([]flatbuffers.UOffsetT, objectsMetaCount)
	i := objectsMetaCount - 1
	for _, meta := range objectsMeta {
		InitBuffers.ObjectMetaStart(builder)
		InitBuffers.ObjectMetaAddObjectType(builder, meta.objectType)
		InitBuffers.ObjectMetaAddCollisionRadius(builder, meta.collisionRadius)
		objectsMetaBuffers[i] = InitBuffers.ObjectMetaEnd(builder)
		i--
	}

	InitBuffers.InitStartObjectsMetaVector(builder, objectsMetaCount)
	for _, buffer := range objectsMetaBuffers {
		builder.PrependUOffsetT(buffer)
	}
	objectsMetaBufObj := builder.EndVector(objectsMetaCount)

	tileLayersCount := len(w.worldmap.TileLayers)
	tileLayersBuffers := make([]flatbuffers.UOffsetT, tileLayersCount)
	for layerIndex, layer := range w.worldmap.TileLayers {
		tileIdsCount := len(layer.TileIds)
		WorldMapBuffers.TileLayerStartTileIdsVector(builder, tileIdsCount)
		for i := range layer.TileIds {
			builder.PrependUint16(layer.TileIds[tileIdsCount-i-1])
		}
		tileIdsBuffObject := builder.EndVector(tileIdsCount)
		WorldMapBuffers.TileLayerStart(builder)
		WorldMapBuffers.TileLayerAddTileIds(builder, tileIdsBuffObject)
		tileLayersBuffers[tileLayersCount-layerIndex-1] = WorldMapBuffers.TileLayerEnd(builder)
	}
	WorldMapBuffers.WorldMapStartLayersVector(builder, tileLayersCount)
	for _, buffer := range tileLayersBuffers {
		builder.PrependUOffsetT(buffer)
	}
	layersBufferObj := builder.EndVector(tileLayersCount)

	WorldMapBuffers.WorldMapStart(builder)
	WorldMapBuffers.WorldMapAddLayers(builder, layersBufferObj)
	WorldMapBuffers.WorldMapAddWidth(builder, int32(w.worldmap.Width))
	WorldMapBuffers.WorldMapAddHeight(builder, int32(w.worldmap.Height))
	worldMapBuffObj := WorldMapBuffers.WorldMapEnd(builder)

	InitBuffers.InitStart(builder)
	InitBuffers.InitAddTimer(builder, timerBufferObj)
	InitBuffers.InitAddWorldMap(builder, worldMapBuffObj)
	InitBuffers.InitAddObjectsMeta(builder, objectsMetaBufObj)
	initBufferObj := InitBuffers.InitEnd(builder)
	builder.Finish(initBufferObj)
	buf := builder.FinishedBytes()

	return append([]byte{byte(CommandsBuffer.CommandInit)}, buf...)
}
