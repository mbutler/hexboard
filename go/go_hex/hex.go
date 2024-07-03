package main

import (
	"encoding/json"
	"fmt"
	"os"
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

func LoadGridFromJSON(filename string) (*Grid, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	jsonSafeDict := make(map[string]map[string]interface{})
	err = json.Unmarshal(data, &jsonSafeDict)
	if err != nil {
		return nil, err
	}

	size := 0
	grid := make(map[Coordinate]map[string]interface{})
	for k, v := range jsonSafeDict {
		var x, y, z int
		fmt.Sscanf(k, "(%d,%d,%d)", &x, &y, &z)
		coords := Coordinate{x, y, z}
		grid[coords] = v
		if abs(x) > size {
			size = abs(x)
		}
		if abs(y) > size {
			size = abs(y)
		}
		if abs(z) > size {
			size = abs(z)
		}
	}

	return &Grid{Size: size, Grid: grid}, nil
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

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
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
		elevation := (p.Noise2D(px, py) + 1) / 2 // Normalize to [0, 1]

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
		{0.0, "deep_ocean"},
		{0.1, "ocean"},
		{0.2, "beach"},
		{0.3, "plains"},
		{0.4, "grassland"},
		{0.5, "forest"},
		{0.6, "shrubland"},
		{0.7, "hills"},
		{0.8, "mountain_foot"},
		{0.9, "mountain"},
		{1.0, "tundra"},
	}

	for _, t := range terrainTypes {
		if elevation < t.threshold {
			return t.terrain
		}
	}
	return "unknown"
}

func main() {
	const gridSize = 1000
	const scale = 0.1

	var grid *Grid
	var err error

	// Try to load the grid from hex_grid.json
	grid, err = LoadGridFromJSON("hex_grid.json")
	if err != nil {
		// If loading fails, generate a new grid
		fmt.Println("Failed to load grid, generating a new one...")
		grid = NewGrid(gridSize)

		// Generate terrain for the grid
		seed := time.Now().UnixNano()
		generateTerrain(grid, scale, seed)

		// Save the new grid to a JSON file
		jsonData, err := grid.ToJSON()
		if err != nil {
			panic(err)
		}
		err = os.WriteFile("hex_grid.json", []byte(jsonData), 0644)
		if err != nil {
			panic(err)
		}
		fmt.Println("Generated and saved new grid to hex_grid.json")
	}

	// Get neighbors of a specific hex in the loaded or newly generated grid
	neighbors := grid.Neighbors(Coordinate{X: 0, Y: 0, Z: 0})
	for _, neighbor := range neighbors {
		fmt.Printf("Neighbor: %+v\n", neighbor)
		fmt.Printf("Properties: Terrain=%s, Elevation=%.2f\n", grid.Grid[neighbor]["terrain"], grid.Grid[neighbor]["elevation"])
	}
}
