class SettlementGenerator {
  constructor() {
    this.settlementRules = [
      { type: 'Single Dwelling', range: [1, 3] },
      { type: 'Thorp', range: [4, 5] },
      { type: 'Hamlet', range: [6, 7] },
      { type: 'Village', range: [8, 9] },
      { type: 'Town', range: [10, 10] },
      { type: 'City', range: [11, 11] },
      { type: 'Castle', range: [12, 14] },
      { type: 'Ruins', range: [15, 16] },
      { type: 'Uninhabited', range: [17, 100] }
    ]
    this.reRollRules = [
      { type: 'Ruined Village', range: [1, 30] },
      { type: 'Ruined City', range: [31, 60] },
      { type: 'Ruined Shrine', range: [61, 85] },
      { type: 'Ruined Tomb', range: [86, 100] }
    ]
    this.castleSizeRules = [
      { size: 'Small', range: [1, 10], type: 'Small shell keep' },
      { size: 'Small', range: [11, 25], type: 'Tower' },
      { size: 'Small', range: [26, 35], type: 'Moat house or friary' },
      { size: 'Medium', range: [36, 45], type: 'Large shell keep' },
      { size: 'Medium', range: [46, 65], type: 'Small walled castle with keep' },
      { size: 'Medium', range: [66, 80], type: 'Medium walled castle with keep' },
      { size: 'Large', range: [81, 88], type: 'Concentric castle' },
      { size: 'Large', range: [89, 95], type: 'Large walled castle with keep' },
      { size: 'Large', range: [96, 100], type: 'Fortress complex' }
    ]
    this.castleInhabitantsRules = {
      'Small': [
        { type: 'Totally deserted', range: [1, 45] },
        { type: 'Deserted (monster therein)', range: [46, 60] },
        { type: 'Humans', range: [61, 70] },
        { type: 'Character-types', range: [71, 100] }
      ],
      'Medium': [
        { type: 'Totally deserted', range: [1, 30] },
        { type: 'Deserted (monster therein)', range: [31, 50] },
        { type: 'Humans', range: [51, 65] },
        { type: 'Character-types', range: [66, 100] }
      ],
      'Large': [
        { type: 'Totally deserted', range: [1, 15] },
        { type: 'Deserted (monster therein)', range: [16, 40] },
        { type: 'Humans', range: [41, 60] },
        { type: 'Character-types', range: [61, 100] }
      ]
    }
  }

  getSettlementType(roll) {
    for (const rule of this.settlementRules) {
      if (roll >= rule.range[0] && roll <= rule.range[1]) {
        if (rule.type === 'Ruins') {
          return this.getReRolledSettlementType()
        } else if (rule.type === 'Castle') {
          return this.getCastleDetails()
        }
        return rule.type
      }
    }
    return 'Uninhabited'
  }

  getReRolledSettlementType() {
    const reRoll = Math.floor(Math.random() * 100) + 1
    for (const rule of this.reRollRules) {
      if (reRoll >= rule.range[0] && reRoll <= rule.range[1]) {
        return rule.type
      }
    }
    return 'Uninhabited'
  }

  getCastleDetails() {
    const castleRoll = Math.floor(Math.random() * 100) + 1
    let castleSize = 'Unknown'
    let castleType = 'Unknown'
    
    for (const rule of this.castleSizeRules) {
      if (castleRoll >= rule.range[0] && castleRoll <= rule.range[1]) {
        castleSize = rule.size
        castleType = rule.type
        break
      }
    }

    const inhabitantRoll = Math.floor(Math.random() * 100) + 1
    let inhabitants = 'Unknown'
    
    for (const rule of this.castleInhabitantsRules[castleSize]) {
      if (inhabitantRoll >= rule.range[0] && inhabitantRoll <= rule.range[1]) {
        inhabitants = rule.type
        break
      }
    }

    return `Castle (${castleType}, ${inhabitants})`
  }
}

export default SettlementGenerator
