import { female, male, surname, places} from './names.js'

class CharacterGenerator {
    constructor(gender, lastName) {
        this.gender = gender || (Math.random() < 0.5 ? 'male' : 'female')
        this.lastName = lastName || surname[Math.floor(Math.random() * surname.length)]
        this.firstName = this.gender === 'male' 
          ? male[Math.floor(Math.random() * male.length)] 
          : female[Math.floor(Math.random() * female.length)]
        this.name = `${this.firstName} ${this.lastName}`
        this.stats = this.rollStats()
        this.qualifiedClasses = this.determineClasses()
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
          qualifiedClasses: this.qualifiedClasses
        }
    }
  }



  export default CharacterGenerator