# Conway's Game of Life

`conways-gol` is a *Conway's Game of Life* implementation with Go and OpenGL using [go-gl](https://github.com/go-gl).

![Conway's Game of Life](./demo.gif)

## *The 'Game'*

Wikipedia describes [Conway's Game of Life](https://en.wikipedia.org/wiki/Conway's_Game_of_Life) as *a cellular automaton devised by the British mathematician John Horton Conway in 1970*.

The premise of the game is that each cell on the grid is, at any time, either dead or alive. The state of each cell is determined using the following rules:

1. Any live cell with fewer than two live neighbours dies, as if caused by  underpopulation.
2. Any live cell with two or three live neighbours lives on to the next generation.
3. Any live cell with more than three live neighbours dies, as if by overpopulation.
4. Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.

For the full rules, check [Wikipedia](https://en.wikipedia.org/wiki/Conway's_Game_of_Life#Rules).

## Author

`conways-gol` was developed by [Kyle Banks](https://twitter.com/kylewbanks).

## License

`conways-gol` is available under the [MIT](./LICENSE) license.