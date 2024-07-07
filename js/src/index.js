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

  let currentTerrain = 'plain' // Ensure starting terrain is set to a valid value

  for (const coords in grid.grid) {
    const terrainRoll = Math.floor(Math.random() * 20) + 1
    const settlementRoll = Math.floor(Math.random() * 100) + 1

    try {
      const terrainResult = terrainGenerator.getTerrainType(currentTerrain, terrainRoll)
      currentTerrain = terrainResult.terrainType
      const notes = terrainResult.notes

      const settlementResult = settlementGenerator.getSettlementType(settlementRoll)
      const { type: settlementType, population } = settlementResult
      const [x, y, z] = coords.split(',').map(Number)
      grid.setProperties({ x, y, z }, { terrain: currentTerrain, settlement: settlementType, population, notes })
    } catch (error) {
      console.error(error.message)
      continue
    }
  }

  const jsonData = grid.toJSON()
  fs.writeFileSync('hex_grid.json', jsonData)

  const elapsedTime = (Date.now() - startTime) / 1000
  console.log(`Elapsed time: ${elapsedTime.toFixed(2)} seconds`)
}

main()
