import fs from 'fs'
import Grid from './grid.js'
import TerrainGenerator from './terrain.js'
import SettlementGenerator from './settlement.js'

async function main() {
  const startTime = Date.now()

  const gridSize = 5
  const grid = new Grid(gridSize)
  const terrainGenerator = new TerrainGenerator()
  const settlementGenerator = new SettlementGenerator()

  // Assuming we start from 'plain' terrain
  let currentTerrain = 'plain'

  // Iterate through each hex and assign terrain based on transition rules
  for (const coords in grid.grid) {
    const terrainRoll = Math.floor(Math.random() * 20) + 1
    const settlementRoll = Math.floor(Math.random() * 100) + 1

    try {
      currentTerrain = terrainGenerator.getTerrainType(currentTerrain, terrainRoll)
    } catch (error) {
      console.error(error.message)
      continue
    }

    const settlementType = settlementGenerator.getSettlementType(settlementRoll)
    const [x, y, z] = coords.split(',').map(Number)
    grid.setProperties({ x, y, z }, { terrain: currentTerrain, settlement: settlementType })
  }

  // Optional: Save the grid with terrain and settlements to a JSON file for inspection
  const jsonData = grid.toJSON()
  fs.writeFileSync('hex_grid_with_terrain_and_settlements.json', jsonData)

  const elapsedTime = (Date.now() - startTime) / 1000
  console.log(`Elapsed time: ${elapsedTime.toFixed(2)} seconds`)
}

main()
