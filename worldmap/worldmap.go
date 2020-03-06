package worldmap

import (
	"aakimov/marsgame/tmx"

	"fmt"
	"log"
)

type WorldMap struct {
	TileLayers    []*TileLayer
	MaterialLayer []uint8
	Width         int
	Height        int
}

func NewWorldMap() *WorldMap {
	return &WorldMap{}
}

type TileLayer struct {
	Name    string
	TileIds []uint16
}

func (wm *WorldMap) Parse(tmxFilePath string) {
	log.Println("start")

	tiledMap := tmx.DecodeByFilePath(tmxFilePath)

	layersCount := len(tiledMap.Layers)
	if layersCount == 0 {
		log.Fatalln("tmx has no layers")
	}

	wm.Width = tiledMap.Width
	wm.Height = tiledMap.Height

	wm.TileLayers = make([]*TileLayer, layersCount)
	for layerIndex, layer := range tiledMap.Layers {
		refs, err := layer.TileGlobalRefs()
		if err != nil {
			panic(err)
		}

		tileIds := make([]uint16, len(refs))
		for i, v := range refs {
			tileIds[i] = matchInternalTiledIdToAssetTileId(v.GlobalID)
		}

		wm.TileLayers[layerIndex] = &TileLayer{
			Name:    layer.Name,
			TileIds: tileIds,
		}
	}
}

var tileMatchingMap = map[tmx.GlobalID]uint16{
	0:   0,
	1:   2,
	2:   3,
	21:  22,
	83:  84,
	84:  85,
	85:  86,
	89:  90,
	90:  91,
	91:  92,
	115: 116,
	116: 117,
	117: 118,
	121: 122,
	122: 123,
	123: 124,
	147: 148,
	148: 149,
	149: 150,
	153: 154,
	154: 155,
	155: 156,
	179: 180,
	180: 181,
	181: 182,
	466: 467,
	497: 498,
	513: 514,
	528: 529,
	529: 530,
	530: 531,
	560: 561,
	561: 562,
	562: 563,
	592: 593,
	593: 594,
	594: 595,
	624: 625,
	625: 626,
}

func matchInternalTiledIdToAssetTileId(gid tmx.GlobalID) uint16 {
	tileId, ok := tileMatchingMap[gid]
	if !ok {
		log.Printf("Lack of tile mathing for GID '%v'", gid)
		return 9999
	}
	return tileId
}

func (wm *WorldMap) PrintToConsole() {
	for _, layer := range wm.TileLayers {
		fmt.Printf("\n%s\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n", layer.Name)
		for i, tileId := range layer.TileIds {
			fmt.Printf("%5d", tileId)
			if (i+1)%wm.Width == 0 {
				fmt.Println("\\")
			}
		}
	}
	fmt.Println()
	log.Fatalln("end for now =)")
}
