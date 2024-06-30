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
	"runtime/pprof"
	"time"
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

func (g *Grid) DrawGrid(size int, outputFile string) {
	minX, minY, maxX, maxY := math.MaxFloat64, math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64
	// First pass: determine image boundaries
	for coords := range g.Grid {
		pixelCoords := g.HexToPixel(coords, float64(size), "flat")
		minX = math.Min(minX, float64(pixelCoords.X))
		minY = math.Min(minY, float64(pixelCoords.Y))
		maxX = math.Max(maxX, float64(pixelCoords.X))
		maxY = math.Max(maxY, float64(pixelCoords.Y))
	}

	margin := float64(size)
	imgWidth := int(maxX - minX + 2*margin)
	imgHeight := int(maxY - minY + 2*margin)
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 255}}, image.Point{}, draw.Src)

	// Precompute hex corners
	hexCorners := make([][]image.Point, len(g.Grid))
	i := 0
	for coords := range g.Grid {
		pixelCoords := g.HexToPixel(coords, float64(size), "flat")
		pixelCoords.X = int(float64(pixelCoords.X) - minX + margin)
		pixelCoords.Y = int(float64(pixelCoords.Y) - minY + margin)

		corners := make([]image.Point, 6)
		for j := 0; j < 6; j++ {
			corner := g.FlatHexCorner(pixelCoords, float64(size), j)
			corners[j] = image.Point{int(corner.X), int(corner.Y)}
		}
		hexCorners[i] = corners
		i++
	}

	// Draw hexagons
	for _, corners := range hexCorners {
		for j := 0; j < len(corners)-1; j++ {
			drawLine(img, corners[j], corners[j+1], color.White)
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

func main() {
	// Profiling start
	f, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	startTime := time.Now()

	// Create a large grid with size 116 (6 miles hexes, size of U.S.)
	gridSize := 116
	grid := NewGrid(gridSize)

	// Set some properties for a few hexes
	for i := -gridSize; i <= gridSize; i++ {
		for j := -gridSize; j <= gridSize; j++ {
			k := -i - j
			if i+j+k == 0 {
				grid.SetProperties(Coordinate{i, j, k}, map[string]interface{}{"value": i * j * k})
			}
		}
	}

	grid.SetProperties(Coordinate{X: 0, Y: 0, Z: 0}, map[string]interface{}{
		"terrain":   "grassland",
		"elevation": 100,
	})

	// get a list of all neighbors of a hex (-1,-78,79)
	neighbors := grid.Neighbors(Coordinate{-1, -78, 79})
	fmt.Printf("Neighbors of (-1,-78,79): %+v\n", neighbors)

	// Perform some operations to test performance
	center := Coordinate{0, 0, 0}
	rangeSize := 1

	hexesInRange := grid.HexesInRange(center, rangeSize, false)
	fmt.Printf("Number of hexes in range: %d\n", len(hexesInRange))

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

	// Draw the grid to a PNG file. This is expensive and slow
	// grid.DrawGrid(10, "hex_grid.png")

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("Elapsed time: %s\n", elapsedTime)
}
