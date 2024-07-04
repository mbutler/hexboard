class TerrainGenerator {
  constructor() {
    this.transitionRules = {
      'plain':          { ranges: [11, 1, 1, 1, 1, 1, 1, 1, 1, 1], types: ['plain', 'scrub', 'forest', 'rough', 'desert', 'hills', 'mountains', 'marsh', 'pond', 'depression'] },
      'scrub':          { ranges: [3, 11, 1, 1, 1, 1, 1, 1, 1, 1], types: ['plain', 'scrub', 'forest', 'rough', 'desert', 'hills', 'mountains', 'marsh', 'pond', 'depression'] },
      'forest':         { ranges: [1, 4, 14, 1, 1, 1, 1, 1, 1, 1], types: ['plain', 'scrub', 'forest', 'rough', 'desert', 'hills', 'mountains', 'marsh', 'pond', 'depression'] },
      'forested hills': { ranges: [1, 1, 1, 1, 1, 4, 14, 1, 1, 1], types: ['plain', 'scrub', 'forest', 'rough', 'desert', 'hills', 'mountains', 'marsh', 'pond', 'depression'] },
      'rough':          { ranges: [1, 1, 5, 8, 1, 1, 1, 1, 1, 1], types: ['plain', 'scrub', 'forest', 'rough', 'desert', 'hills', 'mountains', 'marsh', 'pond', 'depression'] },
      'desert':         { ranges: [3, 5, 4, 14, 1, 1, 1, 1, 1, 1], types: ['plain', 'scrub', 'forest', 'rough', 'desert', 'hills', 'mountains', 'marsh', 'pond', 'depression'] },
      'hills':          { ranges: [1, 2, 3, 4, 5, 14, 1, 1, 1, 1], types: ['plain', 'scrub', 'forest', 'rough', 'desert', 'hills', 'mountains', 'marsh', 'pond', 'depression'] },
      'hilly forest':   { ranges: [1, 1, 1, 1, 1, 2, 3, 4, 5, 14], types: ['plain', 'scrub', 'forest', 'rough', 'desert', 'hills', 'mountains', 'marsh', 'pond', 'depression'] },  // Same as 'hills
      'mountains':      { ranges: [1, 2, 3, 4, 5, 10, 18, 1, 1, 1], types: ['plain', 'scrub', 'forest', 'rough', 'desert', 'hills', 'mountains', 'marsh', 'pond', 'depression'] },
      'mountain pass':  { ranges: [1, 2, 3, 4, 5, 10, 18, 1, 1, 1], types: ['plain', 'scrub', 'forest', 'rough', 'desert', 'hills', 'mountains', 'marsh', 'pond', 'depression'] },  // Same as 'mountains
      'marsh':          { ranges: [2, 4, 5, 6, 8, 1, 15, 1, 1, 1], types: ['plain', 'scrub', 'forest', 'rough', 'desert', 'hills', 'mountains', 'marsh', 'pond', 'depression'] },
      'pond':           { ranges: [1, 1, 1, 1, 1, 1, 1, 19, 1, 1], types: ['plain', 'scrub', 'forest', 'rough', 'desert', 'hills', 'mountains', 'marsh', 'pond', 'depression'] },
      'depression':     { ranges: [1, 1, 1, 1, 1, 1, 1, 1, 1, 20], types: ['plain', 'scrub', 'forest', 'rough', 'desert', 'hills', 'mountains', 'marsh', 'pond', 'depression'] }
    }
  }

  getTerrainType(startTerrain, diceRoll) {
    const terrainRule = this.transitionRules[startTerrain]
    if (!terrainRule) {
      throw new Error(`Invalid terrain type: ${startTerrain}`)
    }

    const { ranges, types } = terrainRule
    let cumulative = 0
    for (let i = 0; i < ranges.length; i++) {
      cumulative += ranges[i]
      if (diceRoll <= cumulative) {
        let terrainType = types[i]
        if (terrainType === 'forest' && Math.random() < 0.1) {
          return 'forested hills'
        } else if (terrainType === 'hills' && Math.random() < 0.1) {
          return 'hilly forest'
        } else if (terrainType === 'mountains' && Math.random() < 0.05) {
          return 'mountain pass'
        }
        return terrainType
      }
    }
    return 'unknown'  // Fallback, should not happen if rules are well defined
  }
}

export default TerrainGenerator
