class TerrainGenerator {
  constructor() {
    this.transitionRules = {
      'plain': [
        { min: 1, max: 11, type: 'plain' },
        { min: 12, max: 12, type: 'scrub' },
        { min: 13, max: 13, type: 'forest' },
        { min: 14, max: 14, type: 'rough' },
        { min: 15, max: 15, type: 'desert' },
        { min: 16, max: 16, type: 'hills' },
        { min: 17, max: 17, type: 'mountains' },
        { min: 18, max: 18, type: 'marsh' },
        { min: 19, max: 19, type: 'pond' },
        { min: 20, max: 20, type: 'depression' }
      ],
      'scrub': [
        { min: 1, max: 3, type: 'plain' },
        { min: 4, max: 11, type: 'scrub' },
        { min: 12, max: 13, type: 'forest' },
        { min: 14, max: 14, type: 'rough' },
        { min: 15, max: 15, type: 'desert' },
        { min: 16, max: 16, type: 'hills' },
        { min: 17, max: 17, type: 'mountains' },
        { min: 18, max: 18, type: 'marsh' },
        { min: 19, max: 19, type: 'pond' },
        { min: 20, max: 20, type: 'depression' }
      ],
      'forest': [
        { min: 1, max: 1, type: 'plain' },
        { min: 2, max: 4, type: 'scrub' },
        { min: 5, max: 14, type: 'forest' },
        { min: 15, max: 15, type: 'rough' },
        { min: 16, max: 16, type: 'hills' },
        { min: 17, max: 17, type: 'mountains' },
        { min: 18, max: 18, type: 'marsh' },
        { min: 19, max: 19, type: 'pond' },
        { min: 20, max: 20, type: 'depression' }
      ],
      'rough': [
        { min: 1, max: 2, type: 'plain' },
        { min: 3, max: 4, type: 'scrub' },
        { min: 5, max: 5, type: 'forest' },
        { min: 6, max: 8, type: 'rough' },
        { min: 9, max: 10, type: 'desert' },
        { min: 11, max: 15, type: 'hills' },
        { min: 16, max: 17, type: 'mountains' },
        { min: 18, max: 18, type: 'marsh' },
        { min: 19, max: 19, type: 'pond' },
        { min: 20, max: 20, type: 'depression' }
      ],
      'desert': [
        { min: 1, max: 3, type: 'plain' },
        { min: 4, max: 5, type: 'scrub' },
        { min: 6, max: 8, type: 'rough' },
        { min: 9, max: 14, type: 'desert' },
        { min: 15, max: 15, type: 'hills' },
        { min: 16, max: 17, type: 'mountains' },
        { min: 18, max: 18, type: 'marsh' },
        { min: 19, max: 19, type: 'pond' },
        { min: 20, max: 20, type: 'depression' }
      ],
      'hills': [
        { min: 1, max: 1, type: 'plain' },
        { min: 2, max: 3, type: 'scrub' },
        { min: 4, max: 5, type: 'forest' },
        { min: 6, max: 7, type: 'rough' },
        { min: 8, max: 8, type: 'desert' },
        { min: 9, max: 14, type: 'hills' },
        { min: 15, max: 16, type: 'mountains' },
        { min: 17, max: 17, type: 'marsh' },
        { min: 18, max: 19, type: 'pond' },
        { min: 20, max: 20, type: 'depression' }
      ],
      'mountains': [
        { min: 1, max: 1, type: 'plain' },
        { min: 2, max: 2, type: 'scrub' },
        { min: 3, max: 3, type: 'forest' },
        { min: 4, max: 5, type: 'rough' },
        { min: 6, max: 6, type: 'desert' },
        { min: 7, max: 10, type: 'hills' },
        { min: 11, max: 18, type: 'mountains' },
        { min: 19, max: 19, type: 'pond' },
        { min: 20, max: 20, type: 'depression' }
      ],
      'marsh': [
        { min: 1, max: 2, type: 'plain' },
        { min: 3, max: 4, type: 'scrub' },
        { min: 5, max: 6, type: 'forest' },
        { min: 7, max: 7, type: 'rough' },
        { min: 8, max: 8, type: 'hills' },
        { min: 9, max: 15, type: 'marsh' },
        { min: 16, max: 16, type: 'mountains' },
        { min: 17, max: 18, type: 'pond' },
        { min: 19, max: 20, type: 'depression' }
      ]
    }
  }

  getTerrainType(startTerrain, diceRoll) {
    const terrainRule = this.transitionRules[startTerrain]
    if (!terrainRule) {
      throw new Error(`Invalid terrain type: ${startTerrain}`)
    }

    let terrainType = startTerrain
    let notes = null

    for (const rule of terrainRule) {
      if (diceRoll >= rule.min && diceRoll <= rule.max) {
        terrainType = rule.type

        // Handle 'pond' and 'depression'
        if (terrainType === 'pond' || terrainType === 'depression') {
          notes = terrainType
          terrainType = startTerrain
        }

        break
      }
    }

    // Handle special cases
    if (terrainType === 'forest' && Math.random() < 0.1) {
      notes = 'forested hills'
    } else if (terrainType === 'hills' && Math.random() < 0.1) {
      notes = 'hilly forest'
    } else if (terrainType === 'mountains' && Math.random() < 0.05) {
      notes = 'mountain pass'
    }

    return { terrainType, notes }
  }
}

export default TerrainGenerator
