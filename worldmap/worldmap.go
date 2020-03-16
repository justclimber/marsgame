package worldmap

import (
	"github.com/justclimber/marsgame/flatbuffers/generated/WalBuffers"
	"github.com/justclimber/marsgame/physics"
	"github.com/justclimber/marsgame/tmx"

	"fmt"
	"log"
)

var EntityTypeMap = map[string]WalBuffers.ObjectType{
	"Player": WalBuffers.ObjectTypeplayer,
	"Xelon":  WalBuffers.ObjectTypexelon,
	"Spore":  WalBuffers.ObjectTypespore,
}

type WorldMap struct {
	TileLayers    []*TileLayer
	MaterialLayer []uint8
	Entities      []Entity
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

type Entity struct {
	Id         uint32
	EntityType WalBuffers.ObjectType
	Pos        physics.Point
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
	if len(tiledMap.ObjectGroups) > 0 {
		wm.Entities = make([]Entity, len(tiledMap.ObjectGroups[0].Objects))
		i := 0
		for _, object := range tiledMap.ObjectGroups[0].Objects {
			objType, ok := EntityTypeMap[object.Type]
			if !ok {
				log.Fatalf("Unsupported entity type %s", object.Type)
			}
			wm.Entities[i] = Entity{
				Id:         uint32(object.ObjectID),
				EntityType: objType,
				Pos:        physics.Point{X: object.X, Y: object.Y},
			}
			i++
		}
	}
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
	0: 0,

	26: 25,
	27: 26,

	58: 54,
	59: 55,

	89: 82,
	90: 83,
	91: 84,

	121: 112,
	122: 113,
	123: 114,

	153: 142,
	154: 143,
	155: 144,

	241: 208,
	242: 209,

	273: 238,
	274: 239,

	304: 267,
	305: 268,
	306: 269,

	336: 297,
	337: 298,
	338: 299,

	368: 327,
	369: 328,
	370: 329,

	400: 356,
	401: 357,
	402: 358,

	770: 356,
	771: 357,
	772: 358,
	773: 358,

	1861: 1372,
	1862: 1373,
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
