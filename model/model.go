package model

import (
	"fmt"
	"math/rand"
	"time"
)

const ( // in pixels
	RadiusSnake	int = 9  // distance from center to side or height of the image
	RadiusFood	int = 5
	NumBlocks	int = 26
	Width		int = (RadiusSnake*2+1)*NumBlocks
	Height		int = (RadiusSnake*2+1)*NumBlocks
)

// direction of the arrows and snake
const (
	LEFT	string = "ArrowLeft"
	RIGHT	string = "ArrowRight"
	UP		string = "ArrowUp"
	DOWN	string = "ArrowDown"
	IDLE	string = "Idle"
)

// SnakeNode will contain the snake as a linked list
type SnakeNode struct {
	X			int
	Y			int
	Next		*SnakeNode
}

// FoodModel contains the coordinates of the current piece of food on the board
type FoodModel struct {
	X			int
	Y			int
}

var (
	Snake			*SnakeNode  // head of the snake
	TailSnake		*SnakeNode  // tail, last node of the snake
	Food			*FoodModel
	ChanKeyPress	chan string  // will contain the 3 most recent key press
	ChanNewGame		chan bool  // will be filled if a new game is requested in the UI
	Score			int
	MaxScore		int // max score in the session
	Direction		string  // the current direction of the snake
	Dead			bool  // true if snake is dead
)

func init() {
	MaxScore = 0  // only initialized once in a session
	NewGame()
}

// NewGame initializes all the model variables to their default values for the start of the game
func NewGame() {
	Snake = &SnakeNode{Width/2, Height/2,  nil}
	TailSnake = Snake  // same as head when length is one
	Direction = IDLE
	rand.Seed(time.Now().UnixNano())
	Food = &FoodModel{}
	GenerateFood()
	ChanKeyPress = make(chan string, 3)
	ChanNewGame = make(chan bool, 1)
	Score = 0
	Dead = false
}

// GenerateFood updates the coordinates of the food to a point that doesn't overlap with the snake
func GenerateFood() {
	// we must check that the food doesnt overlap with the snake
	var x, y int
	var snake *SnakeNode
	// loop creates x, y until it doesnt overlap with the snake
	for {
		x = rand.Intn(Width-2*RadiusFood) + RadiusFood
		y = rand.Intn(Height-2*RadiusFood) + RadiusFood
		overlaps := false
		snake = Snake
		// loop goes through all the nodes of the snake
		for {
			// check for overlap
			if (snake.X-RadiusSnake) <= x && x <= (snake.X+RadiusSnake) &&
				(snake.Y-RadiusSnake) <= y && y <= (snake.Y+RadiusSnake){
				overlaps = true
				break
			}
			// got to the next node of the snake
			snake = snake.Next
			if snake == nil {
				break
			}
		}
		// we break if we found a non-overlapping match
		if !overlaps {
			break
		}
		fmt.Println("In case it gets stuck")
	}
	Food.X = x
	Food.Y = y
}

// PrintSnake prints every node in the snake in order from the head to the tail
func PrintSnake() {
	snake := Snake
	i := 0
	for snake != nil {
		fmt.Printf("Node %d at X:%d, Y:%d\n", i, snake.X, snake.Y)
		snake = snake.Next
	}
	fmt.Println()
}