import fs from 'fs'
import Grid from './grid.js'
import TerrainGenerator from './terrain.js'
import SettlementGenerator from './settlement.js'
import CharacterGenerator from './character.js'
import { places, surname } from './names.js'

function randn_bm() {
  let u = 0, v = 0
  while (u === 0) u = Math.random() // Converting [0,1) to (0,1)
  while (v === 0) v = Math.random()
  return Math.sqrt(-2.0 * Math.log(u)) * Math.cos(2.0 * Math.PI * v)
}

function calculateFamilyDistribution(population) {
  const familyNames = Math.floor(population / 23.08)
  const familyDistribution = new Array(familyNames).fill(0)
  let remainingPopulation = population

  // Generate a normal distribution for family members
  for (let i = 0; i < familyNames; i++) {
    const distribution = Math.abs(randn_bm())
    familyDistribution[i] = distribution
  }

  // Normalize the distribution so it sums to the population
  const sumDistribution = familyDistribution.reduce((a, b) => a + b, 0)
  for (let i = 0; i < familyNames; i++) {
    familyDistribution[i] = Math.floor((familyDistribution[i] / sumDistribution) * population)
  }

  // Adjust any rounding errors by distributing remaining population
  const totalAssigned = familyDistribution.reduce((a, b) => a + b, 0)
  let difference = population - totalAssigned

  while (difference !== 0) {
    const adjustment = difference > 0 ? 1 : -1
    const index = Math.floor(Math.random() * familyNames)
    familyDistribution[index] += adjustment
    difference -= adjustment
  }

  return familyDistribution
}

async function main() {
  const startTime = Date.now()

  // per side
  const gridSize = 1
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
      let notes = terrainResult.notes

      const settlementResult = settlementGenerator.getSettlementType(settlementRoll)
      const { type: settlementType, population } = settlementResult
      const [x, y, z] = coords.split(',').map(Number)

      // Assign a place name to inhabited areas
      if (settlementType !== 'Uninhabited' && settlementType !== 'Ruins' && settlementType !== 'Castle') {
        const placeName = places[Math.floor(Math.random() * places.length)]
        notes = notes ? `${notes}, ${placeName}` : placeName
      }

      // Generate characters for the hex based on the population
      const characters = []
      if (population > 0) {
        const familyDistribution = calculateFamilyDistribution(population)
        for (let i = 0; i < familyDistribution.length; i++) {
          const familySize = familyDistribution[i]
          const familyName = surname[Math.floor(Math.random() * surname.length)]
          for (let j = 0; j < familySize; j++) {
            const character = new CharacterGenerator(familyName)
            characters.push(character.getCharacterDetails())
          }
        }
      }

      grid.setProperties({ x, y, z }, { terrain: currentTerrain, settlement: settlementType, population, notes, characters })
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
  // 31.8 is the area of a hexagon with a size of 6 miles
  // 1122.46 is the area of a hexagon with a size of 36 miles
  console.log(`Total land mass: ${(totalHexes * 31.8).toLocaleString()} square miles`)
  const totalPopulation = Object.values(grid.grid).reduce((total, hex) => total + (hex.population || 0), 0)
  console.log(`Total population: ${totalPopulation.toLocaleString()}`)
  // calculate population density
  const totalLandMass = totalHexes * 31.8
  const populationDensity = totalPopulation / totalLandMass
  console.log(`Population density: ${populationDensity.toFixed(2)} people per square mile`)

  const characterCounts = {}
  let psionicCount = 0

  for (const hex of Object.values(grid.grid)) {
    const { characters } = hex
    if (characters) {
      for (const character of characters) {
        const { qualifiedClasses, psionic } = character

        // classes is an array, count each class and add to running total of characterCounts
        for (const charClass of qualifiedClasses) {
          if (characterCounts[charClass]) {
            characterCounts[charClass]++
          } else {
            characterCounts[charClass] = 1
          }
        }

        if (psionic.isPsionic) {
          psionicCount++
        }
      }
    }
  }

  console.log(characterCounts)
  console.log(`Total number of psionic characters: ${psionicCount}`)
}

main()
