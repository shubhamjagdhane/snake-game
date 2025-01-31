package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
)

type Food struct {
	x, y int
	item byte
}

func newFood(maxRowSize, maxColSize int) *Food {
	return &Food{
		item: '*',
		// TODO: find the exact logic
		x: rand.Intn(maxRowSize-2) + 1,
		y: rand.Intn(maxColSize-2) + 1,
	}
}

// assuming snake can only move from left to right
type Snake struct {
	x, y int
}

type BoardItem byte

const (
	EmptySpace       BoardItem = ' '
	VerticalBorder   BoardItem = '|'
	HorizontalBorder BoardItem = '-'
)

type Board struct {
	height int
	width  int
	data   [][]byte
}

func newBoard(height, width int) *Board {
	return &Board{
		height: height,
		width:  width,
		data:   getBoard(height, width),
	}
}

type Game struct {
	height  int
	width   int
	food    *Food
	snake   *Snake
	board   *Board
	drawBuf *bytes.Buffer
}

func newGame(height, width int) *Game {
	return &Game{
		height:  height,
		width:   width,
		food:    newFood(height, width),
		board:   newBoard(height, width),
		drawBuf: new(bytes.Buffer),
		//		snake:  newSnake(),
	}
}

func (g *Game) RenderBoard() {
	for h := 0; h < g.board.height; h++ {
		for w := 0; w < g.board.width; w++ {
			g.drawBuf.WriteByte(g.board.data[h][w])
		}
		g.drawBuf.WriteByte('\n')
	}
}

func (g *Game) Render() {
	fmt.Fprintln(os.Stdout, g.drawBuf.String())
}

func getBoard(height, width int) [][]byte {
	board := make([][]byte, height)

	for row := range board {
		board[row] = make([]byte, width)

		for col := range board[row] {
			if col == 0 || col == len(board[0])-1 {
				board[row][col] = byte(VerticalBorder)
			} else if row == 0 || row == len(board)-1 {
				board[row][col] = byte(HorizontalBorder)
			} else {
				board[row][col] = byte(EmptySpace)
			}
		}
	}

	return board
}

func main() {
	g := newGame(15, 35)
	g.RenderBoard()
	g.Render()
}
