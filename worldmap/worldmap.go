package worldmap

import (
	"github.com/justclimber/marsgame/flatbuffers/generated/WalBuffers"
	"github.com/justclimber/marsgame/physics"
	"github.com/justclimber/marsgame/tmx"
	"strconv"

	"fmt"
	"log"
)

const TileSize = 32

var EntityTypeMap = map[string]WalBuffers.ObjectType{
	"Player": WalBuffers.ObjectTypeplayer,
	"Xelon":  WalBuffers.ObjectTypexelon,
	"Spore":  WalBuffers.ObjectTypespore,
}

type WorldMap struct {
	TileLayers       []*TileLayer
	MaterialLayer    [][]uint8
	Entities         []Entity
	Width            int
	Height           int
	tileCollisionMap map[uint16]uint8
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

	wm.makeTileCollisionMap(tiledMap)

	layersCount := len(tiledMap.Layers)
	if layersCount == 0 {
		log.Fatalln("tmx has no layers")
	}

	wm.Width = tiledMap.Width
	wm.Height = tiledMap.Height
	wm.MaterialLayer = make([][]uint8, wm.Height)
	for i := range wm.MaterialLayer {
		wm.MaterialLayer[i] = make([]uint8, wm.Width)
	}

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
			id := matchInternalTiledIdToAssetTileId(v.GlobalID)
			tileIds[i] = id
			wm.setCollision(id, i)
		}

		wm.TileLayers[layerIndex] = &TileLayer{
			Name:    layer.Name,
			TileIds: tileIds,
		}
	}
}

func (wm *WorldMap) setCollision(id uint16, i int) {
	collision := wm.tileCollisionMap[id]
	x := i % wm.Width
	y := i / wm.Width
	if wm.MaterialLayer[x][y] < collision {
		wm.MaterialLayer[x][y] = collision
	}
}

func (wm *WorldMap) makeTileCollisionMap(tiledMap *tmx.Map) {
	wm.tileCollisionMap = make(map[uint16]uint8)
	for _, tile := range tiledMap.MainTileSet.Tiles {
		collision := 0
		for _, property := range tile.Properties {
			if property.Name == "collision" {
				var err error
				collision, err = strconv.Atoi(property.Value)
				if err != nil {
					panic(err.Error())
				}
				break
			}
		}
		wm.tileCollisionMap[uint16(tile.TileID)] = uint8(collision)
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

	1477: 10000,
	1478: 10001,
	1479: 10002,
	1480: 10003,
	1481: 10004,
	1482: 10005,
	1483: 10006,
	1484: 10007,
	1485: 10008,

	1509: 10009,
	1510: 10010,
	1511: 10011,
	1512: 10012,
	1513: 10013,
	1514: 10014,
	1515: 10015,
	1516: 10016,
	1517: 10017,

	1573: 10021,
	1574: 10022,
	1575: 10023,
	1576: 10024,
	1577: 10025,
	1578: 10026,
	1579: 10027,
	1580: 10028,
	1581: 10029,

	1607: 10030,
	1608: 10031,
	1609: 10032,
	1610: 10033,
	1611: 10034,
	1612: 10035,
	1613: 10036,
	1614: 10037,
	1615: 10038,

	1547: 10018,
	1548: 10019,
	1549: 10020,

	1643: 10018,
	1644: 10019,
	1645: 10020,
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
