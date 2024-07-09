class Grid {
  constructor(size) {
    this.size = size
    this.grid = this.createHexGrid()
  }

  lerp(a, b, t) {
    return a + (b - a) * t
  }

  movementTable(orientation = 'flat') {
    const directions = orientation === 'flat'
      ? [[1, -1, 0], [0, -1, 1], [-1, 0, 1], [-1, 1, 0], [0, 1, -1], [1, 0, -1]]
      : [[1, 0, -1], [1, -1, 0], [0, -1, 1], [-1, 0, 1], [-1, 1, 0], [0, 1, -1]]

    return directions.reduce((acc, direction, i) => {
      acc[i] = direction
      return acc
    }, {})
  }

  createHexGrid() {
    const grid = {}
    for (let x = -this.size; x <= this.size; x++) {
      for (let y = Math.max(-this.size, -x - this.size); y <= Math.min(this.size, -x + this.size); y++) {
        const z = -x - y
        const coords = `${x},${y},${z}`
        grid[coords] = {}
      }
    }
    return grid
  }

  getHexagon(coords) {
    const { x, y, z } = coords
    return this.grid[`${x},${y},${z}`] || null
  }

  hexesInRange(centerCoords, N, excludeCenter = false) {
    const results = []
    for (let dx = -N; dx <= N; dx++) {
      for (let dy = Math.max(-N, -dx - N); dy <= Math.min(N, -dx + N); dy++) {
        const dz = -dx - dy
        const coords = { x: centerCoords.x + dx, y: centerCoords.y + dy, z: centerCoords.z + dz }
        if (excludeCenter && coords === centerCoords) continue
        results.push(coords)
      }
    }
    return results
  }

  neighbors(coords) {
    return this.hexesInRange(coords, 1, true)
  }

  getProperties(coordinates, prop = null) {
    const coords = `${coordinates.x},${coordinates.y},${coordinates.z}`
    const hexagon = this.grid[coords] || null
    if (!hexagon) return null

    if (prop) {
      return prop.reduce((acc, key) => {
        acc[key] = hexagon[key] || null
        return acc
      }, {})
    } else {
      return hexagon
    }
  }

  setProperties(coordinates, prop) {
    const coords = `${coordinates.x},${coordinates.y},${coordinates.z}`
    if (this.grid[coords]) {
      Object.assign(this.grid[coords], prop)
      return true
    } else {
      return false
    }
  }

  getRelativeCoordinates(startCoords, direction, N) {
    const [dq, dr, ds] = this.movementTable()[direction]
    return {
      x: startCoords.x + dq * N,
      y: startCoords.y + dr * N,
      z: startCoords.z + ds * N
    }
  }

  cubeDistance(startCoords, endCoords) {
    return Math.max(
      Math.abs(startCoords.x - endCoords.x),
      Math.abs(startCoords.y - endCoords.y),
      Math.abs(startCoords.z - endCoords.z)
    )
  }

  cubeRound(cube) {
    let rx = Math.round(cube.x)
    let ry = Math.round(cube.y)
    let rz = Math.round(cube.z)

    const xDiff = Math.abs(rx - cube.x)
    const yDiff = Math.abs(ry - cube.y)
    const zDiff = Math.abs(rz - cube.z)

    if (xDiff > yDiff && xDiff > zDiff) {
      rx = -ry - rz
    } else if (yDiff > zDiff) {
      ry = -rx - rz
    } else {
      rz = -rx - ry
    }

    return { x: rx, y: ry, z: rz }
  }

  hexesInPath(startCoords, endCoords) {
    const N = this.cubeDistance(startCoords, endCoords)
    const results = []
    for (let i = 0; i <= N; i++) {
      const t = 1.0 / N * i
      const cube = {
        x: this.lerp(startCoords.x, endCoords.x, t),
        y: this.lerp(startCoords.y, endCoords.y, t),
        z: this.lerp(startCoords.z, endCoords.z, t)
      }
      results.push(this.cubeRound(cube))
    }
    return results
  }

  hexRangeIntersection(centerA, rangeA, centerB, rangeB) {
    const aList = this.hexesInRange(centerA, rangeA)
    const bList = this.hexesInRange(centerB, rangeB)
    return aList.filter(hexagon => bList.includes(hexagon))
  }

  floodFill(centerCoords, N) {
    const fringes = Array.from({ length: N + 1 }, () => [])
    fringes[0].push(centerCoords)

    for (let k = 1; k <= N; k++) {
      for (const hexCoords of fringes[k - 1]) {
        for (let direction = 0; direction < 6; direction++) {
          const neighborCoords = this.getRelativeCoordinates(hexCoords, direction, 1)
          if (fringes.some(fringe => fringe.includes(neighborCoords))) continue
          const hexProperties = this.getProperties(neighborCoords)
          if (hexProperties?.obstacle) continue
          fringes[k].push(neighborCoords)
        }
      }
    }
    return fringes
  }

  updatePropertiesFromList(hexList) {
    for (const hexObj of hexList) {
      if (hexObj.coords && hexObj.props) {
        const coords = Object.values(hexObj.coords)
        this.setProperties(coords, hexObj.props)
      } else {
        throw new Error("Hexagon object must contain 'coords' and 'props' keys.")
      }
    }
    return this.grid
  }

  directionToIndex(directionStr, orientation = 'flat') {
    const directionMapping = orientation === 'flat'
      ? { 'NE': 0, 'SE': 1, 'S': 2, 'SW': 3, 'NW': 4, 'N': 5 }
      : { 'NE': 0, 'E': 1, 'SE': 2, 'SW': 3, 'W': 4, 'NW': 5 }

    if (!directionMapping[directionStr]) {
      throw new Error(`Invalid direction. Choose from ${Object.keys(directionMapping)}.`)
    }

    return directionMapping[directionStr]
  }

  indexToDirection(directionIndex, orientation = 'flat') {
    const indexMapping = orientation === 'flat'
      ? { 0: 'NE', 1: 'SE', 2: 'S', 3: 'SW', 4: 'NW', 5: 'N' }
      : { 0: 'NE', 1: 'E', 2: 'SE', 3: 'SW', 4: 'W', 5: 'NW' }

    if (!indexMapping[directionIndex]) {
      throw new Error(`Invalid index. Choose from ${Object.keys(indexMapping)}.`)
    }

    return indexMapping[directionIndex]
  }

  flatHexCorner(center, size, i) {
    const angleDeg = 60 * i
    const angleRad = Math.PI / 180 * angleDeg
    return {
      x: center.x + size * Math.cos(angleRad),
      y: center.y + size * Math.sin(angleRad)
    }
  }

  hexToPixel(coords, size, gridSize, orientation = 'flat') {
    const [x, y, z] = coords.split(',').map(Number)
    const px = orientation === 'flat'
      ? size * (3 / 2 * x)
      : size * (Math.sqrt(3) * (x + y / 2))
    const py = orientation === 'flat'
      ? size * (Math.sqrt(3) * (y + x / 2))
      : size * (3 / 2 * y)

    const centerX = size * gridSize
    const centerY = size * gridSize

    return { x: px + centerX, y: py + centerY }
  }

  toJSON() {
    return JSON.stringify(this.grid, null, 4)
  }

  fromJSON(jsonStr) {
    this.grid = JSON.parse(jsonStr)
  }
}

export default Grid
