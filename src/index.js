const SVG = require('svg.js')
import * as Honeycomb from 'honeycomb-grid'
import * as Anime from 'animejs'
import * as Config from './config'

const draw = SVG(Config.map.divContainer)

const Hex = Honeycomb.extendHex({
    size: Config.map.hexSize,
    orientation: Config.map.orientation
})

const Grid = Honeycomb.defineGrid(Hex)

const corners = Hex().corners()

const hexSymbol = draw.symbol()
    .polygon(corners.map(({ x, y }) => `${x},${y}`))
    .fill('none')
    .stroke({ width: Config.map.strokeWidth, color: Config.map.strokeColor })

Grid.rectangle({ width: Config.map.width, height: Config.map.height }).forEach(hex => {
    const { x, y } = hex.toPoint()
    draw.use(hexSymbol).translate(x, y)
})
