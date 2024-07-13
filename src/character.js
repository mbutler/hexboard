import { female, male, surname, places } from './names.js'

class CharacterGenerator {
    constructor(lastName, gender) {
        this.gender = gender || (Math.random() < 0.5 ? 'male' : 'female')
        this.lastName = lastName || surname[Math.floor(Math.random() * surname.length)]
        this.firstName = this.gender === 'male' 
          ? male[Math.floor(Math.random() * male.length)] 
          : female[Math.floor(Math.random() * female.length)]
        this.name = `${this.firstName} ${this.lastName}`
        this.stats = this.rollStats()
        this.qualifiedClasses = this.determineClasses()
        this.psionic = this.determinePsionic()
    }           
  
    // Method to generate a random number between 1 and 6
    rollDice() {
      const roll = () => Math.floor(Math.random() * 6) + 1
      let rolls = [roll(), roll(), roll(), roll()]
      rolls.sort((a, b) => b - a)
      let sum = rolls[0] + rolls[1] + rolls[2]
      return sum
    }
  
    // Method to generate stats for a character
    rollStats() {
      return {
        strength: this.rollDice(),
        intelligence: this.rollDice(),
        wisdom: this.rollDice(),
        dexterity: this.rollDice(),
        constitution: this.rollDice(),
        charisma: this.rollDice()
      }
    }

    // Method to determine psionic ability based on mental stats
    determinePsionic() {
        const { intelligence, wisdom, charisma } = this.stats
        const psionic = { isPsionic: false, strength: 0, ability: 0 }

        if (intelligence < 16 && wisdom < 16 && charisma < 16) {
            return psionic
        }

        let psionicChance = 0
        if (intelligence >= 16 || wisdom >= 16 || charisma >= 16) {
            psionicChance = 1
        }
        if (intelligence > 16) {
            psionicChance += (intelligence - 16) * 2.5
        }
        if (wisdom > 16) {
            psionicChance += (wisdom - 16) * 1.5
        }
        if (charisma > 16) {
            psionicChance += (charisma - 16) * 1.5
        }

        const randomRoll = Math.floor(Math.random() * 100) + 1
        if (randomRoll <= psionicChance) {
            psionic.isPsionic = true
            let mentalStrength = Math.max(intelligence - 12, 0) + Math.max(wisdom - 12, 0) + Math.max(charisma - 12, 0)
            if (intelligence > 16 && wisdom > 16 && charisma > 16) {
                mentalStrength *= 4
            } else if ((intelligence > 16 && wisdom > 16) || (intelligence > 16 && charisma > 16) || (wisdom > 16 && charisma > 16)) {
                mentalStrength *= 2
            }
            psionic.strength = mentalStrength + Math.floor(Math.random() * 100) + 1
            psionic.ability = psionic.strength * 2
        }

        return psionic
    }
  
    // Method to determine qualified classes based on stats
    determineClasses() {
      const classes = []
  
      const requirements = {
        Cleric: {strength: 6, intelligence: 6, wisdom: 9, constitution: 6, charisma: 6},
        Druid: {strength: 6, intelligence: 6, wisdom: 12, dexterity: 6, constitution: 6, charisma: 15},
        Fighter: {strength: 9, wisdom: 6, dexterity: 6, constitution: 7, charisma: 6},
        Paladin: {strength: 12, intelligence: 9, wisdom: 13, dexterity: 6, constitution: 9, charisma: 17},
        Ranger: {strength: 13, intelligence: 13, wisdom: 14, dexterity: 6, constitution: 14, charisma: 6},
        MagicUser: {intelligence: 9, wisdom: 6, dexterity: 6, constitution: 6, charisma: 6},
        Illusionist: {strength: 6, intelligence: 15, wisdom: 6, dexterity: 16, charisma: 6},
        Thief: {strength: 6, intelligence: 6, dexterity: 9, constitution: 6, charisma: 6},
        Assassin: {strength: 12, intelligence: 11, wisdom: 6, dexterity: 12, constitution: 6},
        Monk: {strength: 15, intelligence: 6, wisdom: 15, dexterity: 15, constitution: 11, charisma: 6},
        Bard: {strength: 15, intelligence: 12, wisdom: 15, dexterity: 15, constitution: 10, charisma: 15}
      }
  
      for (let className in requirements) {
        const requirement = requirements[className]
        const qualified = Object.keys(requirement).every(
          key => this.stats[key] >= requirement[key]
        )
        if (qualified) {
          classes.push(className)
        }
      }
  
      return classes
    }

    getCharacterDetails() {
        return {
          gender: this.gender,
          name: this.name,
          stats: this.stats,
          psionic: this.psionic,
          qualifiedClasses: this.qualifiedClasses
        }
    }
}

export default CharacterGenerator
