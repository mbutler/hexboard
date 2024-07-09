import fs from 'fs'
import Grid from './grid.js'
import TerrainGenerator from './terrain.js'
import SettlementGenerator from './settlement.js'
import CharacterGenerator from './character.js'

async function main() {
  const startTime = Date.now()

  const gridSize = 20
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

  const totalHexes = Object.keys(grid.grid).length

  const elapsedTime = (Date.now() - startTime) / 1000
  console.log(`Elapsed time: ${elapsedTime.toFixed(2)} seconds`)
  console.log(`Total number of hexes: ${totalHexes}`)
  // area of a hexagon = 1.5 * sqrt(3) * s^2
  // 1.95 is the area of a hexagon with a size of 1 mile
  console.log(`Total land mass: ${(totalHexes * 1.95).toLocaleString()} square miles`)
  const totalPopulation = Object.values(grid.grid).reduce((total, hex) => total + (hex.population || 0), 0)
  console.log(`Total population: ${totalPopulation.toLocaleString()}`)
  // calculate population density
  const totalLandMass = totalHexes * 1.95
  const populationDensity = totalPopulation / totalLandMass
  console.log(`Population density: ${populationDensity.toFixed(2)} people per square mile`)

  const characterCounts = {}

for (let i = 0; i < totalPopulation; i++) {
  const character = new CharacterGenerator()
  const details = character.getCharacterDetails()
  const classes = details.qualifiedClasses

  // classes is an array, count each class and add to running total of characterCounts
  for (const charClass of classes) {
    if (characterCounts[charClass]) {
      characterCounts[charClass]++
    } else {
      characterCounts[charClass] = 1
    }
  }
}

console.log(characterCounts)

}

main()
