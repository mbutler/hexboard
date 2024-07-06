import { describe, it, expect } from 'bun:test'
import TerrainGenerator from '../src/terrain.js'

describe('TerrainGenerator', () => {
  const generator = new TerrainGenerator()

  it('should return valid terrain types based on dice roll', () => {
    const testCases = [
      { start: 'plain', roll: 1, expected: 'plain' },
      { start: 'plain', roll: 12, expected: 'scrub' },
      { start: 'forest', roll: 5, expected: 'scrub' },
      { start: 'desert', roll: 20, expected: 'rough' },
      { start: 'marsh', roll: 4, expected: 'scrub' }
    ]

    testCases.forEach(({ start, roll, expected }) => {
      expect(generator.getTerrainType(start, roll)).toBe(expected)
    })
  })

  it('should apply special case for forested hills with 10% chance', () => {
    const mockMath = Object.create(global.Math)
    mockMath.random = () => 0.05
    global.Math = mockMath

    expect(generator.getTerrainType('forest', 5)).toBe('scrub')

    global.Math.random = () => 0.15
    expect(generator.getTerrainType('forest', 5)).toBe('scrub')

    global.Math = Object.create(global.Math)
  })

  it('should apply special case for hilly forest with 10% chance', () => {
    const mockMath = Object.create(global.Math)
    mockMath.random = () => 0.05
    global.Math = mockMath

    expect(generator.getTerrainType('hills', 5)).toBe('forested hills')

    global.Math.random = () => 0.15
    expect(generator.getTerrainType('hills', 5)).toBe('forest')

    global.Math = Object.create(global.Math)
  })

  it('should apply special case for mountain pass with 5% chance', () => {
    const mockMath = Object.create(global.Math)
    mockMath.random = () => 0.02
    global.Math = mockMath

    expect(generator.getTerrainType('mountains', 5)).toBe('forested hills')

    global.Math.random = () => 0.06
    expect(generator.getTerrainType('mountains', 5)).toBe('forested hills')

    global.Math = Object.create(global.Math)
  })

  it('should throw error for invalid terrain type', () => {
    expect(() => generator.getTerrainType('invalid', 1)).toThrow('Invalid terrain type: invalid')
  })
})
