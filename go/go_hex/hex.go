package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"time"
)

type Coordinate struct {
	X, Y, Z int `json:"x,y,z"`
}

type Grid struct {
	Size int
	Grid map[Coordinate]map[string]interface{}
}

func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

func movementTable(orientation string) (map[int][3]int, error) {
	var directions [6][3]int
	switch orientation {
	case "flat":
		directions = [6][3]int{{1, -1, 0}, {0, -1, 1}, {-1, 0, 1}, {-1, 1, 0}, {0, 1, -1}, {1, 0, -1}}
	case "pointy":
		directions = [6][3]int{{1, 0, -1}, {1, -1, 0}, {0, -1, 1}, {-1, 0, 1}, {-1, 1, 0}, {0, 1, -1}}
	default:
		return nil, fmt.Errorf("invalid orientation. Choose either 'flat' or 'pointy'")
	}
	movementTable := make(map[int][3]int)
	for i, direction := range directions {
		movementTable[i] = direction
	}
	return movementTable, nil
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

func (g *Grid) GetHexagon(coords Coordinate) map[string]interface{} {
	return g.Grid[coords]
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

func (g *Grid) Neighbors(coords Coordinate) []Coordinate {
	return g.HexesInRange(coords, 1, true)
}

func (g *Grid) GetProperties(coordinates Coordinate, prop []string) map[string]interface{} {
	if hexagon, ok := g.Grid[coordinates]; ok {
		if prop != nil {
			result := make(map[string]interface{})
			for _, p := range prop {
				result[p] = hexagon[p]
			}
			return result
		}
		return hexagon
	}
	return nil
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

func (g *Grid) GetRelativeCoordinates(startCoords Coordinate, direction, N int) (Coordinate, error) {
	directions, err := movementTable("flat")
	if err != nil {
		return Coordinate{}, err
	}
	dq, dr, ds := directions[direction][0], directions[direction][1], directions[direction][2]
	newCoords := Coordinate{
		startCoords.X + dq*N,
		startCoords.Y + dr*N,
		startCoords.Z + ds*N,
	}
	return newCoords, nil
}

func (g *Grid) CubeDistance(startCoords, endCoords Coordinate) int {
	return maxInt(abs(startCoords.X-endCoords.X), abs(startCoords.Y-endCoords.Y), abs(startCoords.Z-endCoords.Z))
}

func (g *Grid) CubeRound(cube map[string]float64) Coordinate {
	rx := round(cube["x"])
	ry := round(cube["y"])
	rz := round(cube["z"])

	xDiff := math.Abs(float64(rx) - cube["x"])
	yDiff := math.Abs(float64(ry) - cube["y"])
	zDiff := math.Abs(float64(rz) - cube["z"])

	if xDiff > yDiff && xDiff > zDiff {
		rx = -ry - rz
	} else if yDiff > zDiff {
		ry = -rx - rz
	} else {
		rz = -rx - ry
	}
	return Coordinate{rx, ry, rz}
}

func (g *Grid) HexesInPath(startCoords, endCoords Coordinate) []Coordinate {
	N := g.CubeDistance(startCoords, endCoords)
	results := []Coordinate{}
	for i := 0; i <= N; i++ {
		t := 1.0 / float64(N) * float64(i)
		cube := map[string]float64{
			"x": lerp(float64(startCoords.X), float64(endCoords.X), t),
			"y": lerp(float64(startCoords.Y), float64(endCoords.Y), t),
			"z": lerp(float64(startCoords.Z), float64(endCoords.Z), t),
		}
		results = append(results, g.CubeRound(cube))
	}
	return results
}

func (g *Grid) HexRangeIntersection(centerA Coordinate, rangeA int, centerB Coordinate, rangeB int) []Coordinate {
	aList := g.HexesInRange(centerA, rangeA, false)
	bList := g.HexesInRange(centerB, rangeB, false)
	intersection := []Coordinate{}
	for _, a := range aList {
		for _, b := range bList {
			if a == b {
				intersection = append(intersection, a)
			}
		}
	}
	return intersection
}

func (g *Grid) FloodFill(centerCoords Coordinate, N int) [][]Coordinate {
	fringes := make([][]Coordinate, N+1)
	fringes[0] = append(fringes[0], centerCoords)

	for k := 1; k <= N; k++ {
		for _, hexCoords := range fringes[k-1] {
			for direction := 0; direction < 6; direction++ {
				neighborCoords, err := g.GetRelativeCoordinates(hexCoords, direction, 1)
				if err != nil {
					continue
				}
				alreadyInFringe := false
				for _, fringe := range fringes {
					for _, coords := range fringe {
						if neighborCoords == coords {
							alreadyInFringe = true
							break
						}
					}
					if alreadyInFringe {
						break
					}
				}
				if alreadyInFringe {
					continue
				}
				if hexProperties := g.GetProperties(neighborCoords, nil); hexProperties != nil {
					if val, ok := hexProperties["obstacle"]; ok && val.(bool) {
						continue
					}
				}
				fringes[k] = append(fringes[k], neighborCoords)
			}
		}
	}
	return fringes
}

func (g *Grid) UpdatePropertiesFromList(hexList []map[string]interface{}) {
	for _, hexObj := range hexList {
		if coords, ok := hexObj["coords"].(map[string]interface{}); ok {
			x := int(coords["x"].(float64))
			y := int(coords["y"].(float64))
			z := int(coords["z"].(float64))
			props := hexObj["props"].(map[string]interface{})
			g.SetProperties(Coordinate{x, y, z}, props)
		} else {
			panic("Hexagon object must contain 'coords' and 'props' keys.")
		}
	}
}

func (g *Grid) DrawGrid(size int, outputFile string) {
	minX, minY, maxX, maxY := math.MaxFloat64, math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64
	for coords := range g.Grid {
		pixelCoords := g.HexToPixel(coords, float64(size), "flat")
		minX = math.Min(minX, float64(pixelCoords.X))
		minY = math.Min(minY, float64(pixelCoords.Y))
		maxX = math.Max(maxX, float64(pixelCoords.X))
		maxY = math.Max(maxY, float64(pixelCoords.Y))
	}

	margin := float64(size)
	imgSize := image.Rect(0, 0, int(maxX-minX+2*margin), int(maxY-minY+2*margin))
	img := image.NewRGBA(imgSize)
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 255}}, image.Point{}, draw.Src)

	for coords := range g.Grid {
		pixelCoords := g.HexToPixel(coords, float64(size), "flat")
		pixelCoords.X = int(float64(pixelCoords.X) - minX + margin)
		pixelCoords.Y = int(float64(pixelCoords.Y) - minY + margin)

		corners := []image.Point{}
		for i := 0; i < 6; i++ {
			corner := g.FlatHexCorner(pixelCoords, float64(size), i)
			corners = append(corners, image.Point{int(corner.X), int(corner.Y)})
		}
		for i := 0; i < len(corners)-1; i++ {
			drawLine(img, corners[i], corners[i+1], color.White)
		}
		drawLine(img, corners[len(corners)-1], corners[0], color.White)
	}

	outFile, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()
	png.Encode(outFile, img)
}

func drawLine(img *image.RGBA, p1, p2 image.Point, col color.Color) {
	dx := abs(p2.X - p1.X)
	dy := abs(p2.Y - p1.Y)
	sx, sy := 1, 1
	if p1.X >= p2.X {
		sx = -1
	}
	if p1.Y >= p2.Y {
		sy = -1
	}
	err := dx - dy

	for {
		img.Set(p1.X, p1.Y, col)
		if p1.X == p2.X && p1.Y == p2.Y {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			p1.X += sx
		}
		if e2 < dx {
			err += dx
			p1.Y += sy
		}
	}
}

func (g *Grid) HexToPixel(coords Coordinate, size float64, orientation string) image.Point {
	var px, py float64
	switch orientation {
	case "flat":
		px = size * (3.0 / 2.0 * float64(coords.X))
		py = size * (math.Sqrt(3) * (float64(coords.Y) + float64(coords.X)/2.0))
	case "pointy":
		px = size * (math.Sqrt(3) * (float64(coords.X) + float64(coords.Y)/2.0))
		py = size * (3.0 / 2.0 * float64(coords.Y))
	default:
		panic("Invalid orientation. Choose either 'flat' or 'pointy'.")
	}
	centerX := size * 10
	centerY := size * 10

	return image.Point{int(px + centerX), int(py + centerY)}
}

func (g *Grid) FlatHexCorner(center image.Point, size float64, i int) image.Point {
	angleDeg := 60 * float64(i)
	angleRad := math.Pi / 180 * angleDeg
	return image.Point{
		int(float64(center.X) + size*math.Cos(angleRad)),
		int(float64(center.Y) + size*math.Sin(angleRad)),
	}
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

func (g *Grid) FromJSON(jsonStr string) error {
	jsonSafeDict := make(map[string]map[string]interface{})
	if err := json.Unmarshal([]byte(jsonStr), &jsonSafeDict); err != nil {
		return err
	}
	g.Grid = make(map[Coordinate]map[string]interface{})
	for k, v := range jsonSafeDict {
		var x, y, z int
		fmt.Sscanf(k, "(%d,%d,%d)", &x, &y, &z)
		g.Grid[Coordinate{x, y, z}] = v
	}
	return nil
}

func round(value float64) int {
	if value < 0 {
		return int(value - 0.5)
	}
	return int(value + 0.5)
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

func maxInt(a, b, c int) int {
	if a > b && a > c {
		return a
	}
	if b > c {
		return b
	}
	return c
}

func main() {
	startTime := time.Now()

	// Create a large grid with size 1000
	gridSize := 100
	grid := NewGrid(gridSize)

	// Set some properties for a few hexes
	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			for k := 0; k < gridSize; k++ {
				if i+j+k == 0 {
					grid.SetProperties(Coordinate{i, j, k}, map[string]interface{}{"value": i * j * k})
				}
			}
		}
	}

	// Perform some operations to test performance
	center := Coordinate{0, 0, 0}
	rangeSize := 50

	hexesInRange := grid.HexesInRange(center, rangeSize, false)
	fmt.Printf("Number of hexes in range: %d\n", len(hexesInRange))

	neighbors := grid.Neighbors(center)
	fmt.Printf("Number of neighbors: %d\n", len(neighbors))

	relativeCoords, err := grid.GetRelativeCoordinates(center, 2, 5)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Relative coordinates: %+v\n", relativeCoords)

	// Write the grid to a JSON file
	jsonData, err := grid.ToJSON()
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("hex_grid.json", []byte(jsonData), 0644)
	if err != nil {
		panic(err)
	}

	// Draw the grid to a PNG file
	grid.DrawGrid(10, "hex_grid.png")

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("Elapsed time: %s\n", elapsedTime)
}
