package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Food struct {
	x, y       int
	maxRowSize int
	maxColSize int
	item       byte
}

func newFood(maxRowSize, maxColSize int) *Food {
	f := &Food{
		item:       '*',
		maxRowSize: maxRowSize,
		maxColSize: maxColSize,
	}

	// to avoid the initial row & col
	f.PlaceFoodOnBoard()

	return f
}

func (f *Food) PlaceFoodOnBoard() {
	f.x = rand.Intn(f.maxRowSize-2) + 1
	f.y = rand.Intn(f.maxColSize-2) + 1
}

// assuming snake can only move from left to right
type Snake struct {
	x, y                   int
	maxRowSize, maxColSize int
	item                   byte
}

func newSnake(maxRowSize, maxColSize int) *Snake {

	return &Snake{
		x:          rand.Intn(maxRowSize-2) + 1,
		y:          rand.Intn(maxColSize-2) + 1,
		maxRowSize: maxRowSize,
		maxColSize: maxColSize,
		item:       '-',
	}
}

func (s *Snake) PlaceSnakeOnBoard() {
	num := rand.Intn(200)
	if num%2 == 0 {
		s.BottomToTop()
	} else if num%3 == 0 {
		s.LeftToRight()
	} else if num%5 == 0 {
		s.RightToLeft()
	} else {
		s.TopToBottom()
	}
}

func (s *Snake) LeftToRight() {
	if s.y == s.maxColSize-2 {
		s.y = 0
	}
	s.y += 1
	s.item = '-'
}

func (s *Snake) RightToLeft() {
	if s.y == 1 {
		s.y = s.maxColSize - 2
	}
	s.y -= 1
	s.item = '-'
}

func (s *Snake) BottomToTop() {
	if s.x == 1 {
		s.x = s.maxRowSize - 1
	}
	s.x -= 1
	s.item = '|'
}

func (s *Snake) TopToBottom() {
	if s.x == s.maxRowSize-2 {
		s.x = 0
	}
	s.x += 1
	s.item = '|'
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
	height            int
	width             int
	food              *Food
	snake             *Snake
	board             *Board
	drawBuf           *bytes.Buffer
	initialFoodRender bool
}

func newGame(height, width int) *Game {
	return &Game{
		height:            height,
		width:             width,
		food:              newFood(height, width),
		board:             newBoard(height, width),
		drawBuf:           new(bytes.Buffer),
		snake:             newSnake(height, width),
		initialFoodRender: true,
	}
}

func (g *Game) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			g.Render()
			time.Sleep(time.Millisecond * 150)
		}
	}
}

func (g *Game) RenderFood() {
	if g.initialFoodRender {
		g.board.data[g.food.x][g.food.y] = byte(EmptySpace)
		g.food.PlaceFoodOnBoard()
		g.board.data[g.food.x][g.food.y] = g.food.item
		g.initialFoodRender = false
	}

	if g.snake.x == g.food.x && g.snake.y == g.food.y {
		g.board.data[g.food.x][g.food.y] = byte(EmptySpace)
		g.food.PlaceFoodOnBoard()
		g.board.data[g.food.x][g.food.y] = g.food.item
	}
}

func (g *Game) RenderSnake() {

	g.board.data[g.snake.x][g.snake.y] = byte(EmptySpace)
	g.snake.PlaceSnakeOnBoard()
	g.board.data[g.snake.x][g.snake.y] = g.snake.item
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
	g.drawBuf.Reset()
	g.RenderBoard()
	g.RenderSnake()
	g.RenderFood()
	fmt.Fprint(os.Stdout, "\033[2J\033[1;1H")
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
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	g := newGame(15, 35)
	go g.Start(ctx)

	<-sigChan
	cancel()
}
