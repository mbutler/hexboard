package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/aquilax/go-perlin"
)

type Coordinate struct {
	X, Y, Z int
}

type Grid struct {
	Size int
	Grid map[Coordinate]map[string]interface{}
}

func NewGrid(size int) *Grid {
	g := &Grid{
		Size: size,
		Grid: make(map[Coordinate]map[string]interface{}),
	}
	g.createHexGrid()
	return g
}

func (g *Grid) createHexGrid() {
	for x := -g.Size; x <= g.Size; x++ {
		for y := max(-g.Size, -x-g.Size); y <= min(g.Size, -x+g.Size); y++ {
			z := -x - y
			coords := Coordinate{x, y, z}
			g.Grid[coords] = make(map[string]interface{})
		}
	}
}

func (g *Grid) SetProperties(coordinates Coordinate, prop map[string]interface{}) bool {
	if hexagon, ok := g.Grid[coordinates]; ok {
		for k, v := range prop {
			hexagon[k] = v
		}
		return true
	}
	return false
}

func (g *Grid) ToJSON() (string, error) {
	jsonSafeDict := make(map[string]map[string]interface{})
	for k, v := range g.Grid {
		jsonSafeDict[fmt.Sprintf("(%d,%d,%d)", k.X, k.Y, k.Z)] = v
	}
	jsonData, err := json.Marshal(jsonSafeDict)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (g *Grid) Neighbors(coords Coordinate) []Coordinate {
	return g.HexesInRange(coords, 1, true)
}

func (g *Grid) HexesInRange(centerCoords Coordinate, N int, excludeCenter bool) []Coordinate {
	results := []Coordinate{}
	for dx := -N; dx <= N; dx++ {
		for dy := max(-N, -dx-N); dy <= min(N, -dx+N); dy++ {
			dz := -dx - dy
			coords := Coordinate{centerCoords.X + dx, centerCoords.Y + dy, centerCoords.Z + dz}
			if excludeCenter && coords == centerCoords {
				continue
			}
			results = append(results, coords)
		}
	}
	return results
}

func (g *Grid) GetRelativeCoordinates(startCoords Coordinate, direction, N int) (Coordinate, error) {
	directions := map[int][3]int{
		0: {1, -1, 0}, 1: {0, -1, 1}, 2: {-1, 0, 1},
		3: {-1, 1, 0}, 4: {0, 1, -1}, 5: {1, 0, -1},
	}
	dq, dr, ds := directions[direction][0], directions[direction][1], directions[direction][2]
	newCoords := Coordinate{
		startCoords.X + dq*N,
		startCoords.Y + dr*N,
		startCoords.Z + ds*N,
	}
	return newCoords, nil
}

func generateTerrain(grid *Grid, scale float64, seed int64) {
	alpha := 2.0
	beta := 2.0
	n := int32(3)
	p := perlin.NewPerlin(alpha, beta, n, seed)

	for coords := range grid.Grid {
		px := float64(coords.X) * scale
		py := float64(coords.Y) * scale
		elevation := p.Noise2D(px, py)

		terrain := determineTerrain(elevation)
		grid.SetProperties(coords, map[string]interface{}{
			"terrain":   terrain,
			"elevation": elevation,
		})
	}
}

func determineTerrain(elevation float64) string {
	terrainTypes := []struct {
		threshold float64
		terrain   string
	}{
		{-0.4, "deep_ocean"},
		{-0.2, "ocean"},
		{0.0, "beach"},
		{0.05, "plains"},
		{0.1, "grassland"},
		{0.15, "forest"},
		{0.2, "savanna"},
		{0.25, "shrubland"},
		{0.3, "hills"},
		{0.35, "mountain_foot"},
		{0.4, "mountain"},
		{0.45, "tundra"},
		{0.5, "snowy_mountain"},
		{0.55, "glacier"},
		{0.6, "volcano"},
		{0.65, "desert"},
		{0.7, "swamp"},
		{0.75, "rainforest"},
		{0.8, "mangrove"},
		{0.85, "coral_reef"},
	}

	for _, t := range terrainTypes {
		if elevation < t.threshold {
			return t.terrain
		}
	}
	return "unknown"
}

func main() {
	// Profiling start
	f, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	startTime := time.Now()

	gridSize := 116
	grid := NewGrid(gridSize)

	for i := -gridSize; i <= gridSize; i++ {
		for j := -gridSize; j <= gridSize; j++ {
			k := -i - j
			if i+j+k == 0 {
				grid.SetProperties(Coordinate{i, j, k}, map[string]interface{}{"value": i * j * k})
			}
		}
	}

	// Generate terrain for the grid
	scale := 0.1 // Adjust scale for different detail levels
	seed := time.Now().UnixNano()
	generateTerrain(grid, scale, seed)

	// Get neighbors of a specific hex
	neighbors := grid.Neighbors(Coordinate{X: 0, Y: 0, Z: 0})
	for _, neighbor := range neighbors {
		fmt.Printf("Neighbor: %+v\n", neighbor)
		fmt.Printf("Properties: Terrain=%s, Elevation=%.2f\n", grid.Grid[neighbor]["terrain"], grid.Grid[neighbor]["elevation"])
	}

	jsonData, err := grid.ToJSON()
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("hex_grid.json", []byte(jsonData), 0644)
	if err != nil {
		panic(err)
	}

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("Elapsed time: %s\n", elapsedTime)
}
