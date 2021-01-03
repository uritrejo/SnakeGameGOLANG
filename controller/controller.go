package controller

import (
	"fmt"
	"../model"
	"time"
)

const (
	loopPeriod	int = 33  //ms
	speedSnake	int = 80  //ms to move a block
)

var (
	lastRecordedTime	time.Time
	dirChangeInP 		bool = false  // 'direction changed in Period'
		// to limit it to once per period, otherwise, if clicked fast enough, the snake could turn on itself
)

// Start runs the loop of the game, will run until process is stopped
func Start() {
	lastRecordedTime = time.Now()

	for {
		// we check if a key was pressed
		select{
		case newDirection := <- model.ChanKeyPress:
			if !dirChangeInP {
				onKeyPressed(newDirection)
			}
			dirChangeInP = true
		default:
			// no keyPressed, nothing
		}

		// we check if there's a request for a new game
		select{
		case <- model.ChanNewGame:
			newGame()
		default:
			// no new game, nothing
		}

		// every speedSnake ms we move the snake
		if time.Now().After(lastRecordedTime.Add(time.Millisecond * time.Duration(speedSnake))) && !model.Dead{
			step()
			lastRecordedTime = time.Now()
			dirChangeInP = false
		}

		// we sleep
		time.Sleep(time.Millisecond * time.Duration(loopPeriod))
	}
}

// onKeyPressed handles the event of a key press, adjusts snake's direction if necessary
func onKeyPressed(keyPressed string) {
	// check the string to see if it matches a direction
	// if so, change the direction of the snake to that one
	if keyPressed == model.LEFT || keyPressed == model.RIGHT ||
		keyPressed == model.UP || keyPressed == model.DOWN {
		//fmt.Printf("%s pressed\n", keyPressed)

		// the snake can't turn towards itself
		if model.Direction == model.LEFT && keyPressed == model.RIGHT ||
			model.Direction == model.RIGHT && keyPressed == model.LEFT ||
			model.Direction == model.UP && keyPressed == model.DOWN ||
			model.Direction == model.DOWN && keyPressed == model.UP {
			return
		}

		model.Direction = keyPressed
	}
}

// step updates the model to represent the new state of the game
func step() {

	if model.Direction == model.IDLE {
		return
	}

	// in temp we'll store the coordinates of each in order to assign it to the successive node as the snake moves
	previousX, previousY := model.Snake.X, model.Snake.Y

	// we also store the coordinates of the last node as it moves, the new node may take these coordinates
	newNodeX, newNodeY := model.TailSnake.X, model.TailSnake.Y

	// we first update the coordinates of the head according to the direction
	switch model.Direction {
	case model.LEFT:
		model.Snake.X -= 2*model.RadiusSnake
	case model.RIGHT:
		model.Snake.X += 2*model.RadiusSnake
	//the y axis on the image increases towards the bottom
	case model.UP:
		model.Snake.Y -= 2*model.RadiusSnake
	case model.DOWN:
		model.Snake.Y += 2*model.RadiusSnake
	default:
		// nothing
	}

	// now we move each node to the previous position of their predecessor
	snake := model.Snake.Next
	for snake != nil {
		// update current position
		previousX, previousY, snake.X, snake.Y = snake.X, snake.Y, previousX, previousY
		snake = snake.Next
	}

	if died() {
		fmt.Println("You have died!")
		// we'll wait for the user to click New Game
		//newGame()
		model.Direction = model.IDLE
		model.Dead = true
		return
	}

	if ateFood() {
		model.Score += 10

		if model.MaxScore < model.Score {
			model.MaxScore = model.Score
		}

		newNode := &model.SnakeNode{newNodeX, newNodeY, nil}
		model.TailSnake.Next = newNode
		model.TailSnake = model.TailSnake.Next

		// we update the coordinates of the food
		model.GenerateFood()
	}
}

// TODO
// newGame resets the model to create a new game
func newGame() {
	fmt.Println("New Game to be started")
	model.NewGame()
	dirChangeInP = false
}

// died returns true if a single pixel of the snake went outside of the boundaries (death)
func died() bool {
	died := false
	snake := model.Snake
	rS := model.RadiusSnake
	// we check for the boundaries of the board
	if snake.X < rS || snake.X > (model.Width - rS) || snake.Y < rS || snake.Y > (model.Height - rS) {
		died = true
	}

	// TODO we check for the snake intersecting itself
	currNode := snake.Next
	for currNode != nil {
		// we check if the head intersects with any of its pieces
		if !((currNode.X-rS) >= (snake.X+rS) || (snake.X-rS) >= (currNode.X+rS)) &&
			!((currNode.Y-rS) >= (snake.Y+rS) || (snake.Y-rS) >= (currNode.Y+rS)) {
			died = true
			break
		}
		currNode = currNode.Next
	}

	return died
}

// ateFood returns true if a pixel of the snake came into contact with a pixel of the food
func ateFood() bool {
	snake := model.Snake
	food := model.Food
	rS := model.RadiusSnake
	rF := model.RadiusFood

	// we check if they align on the X axis
	if (food.X-rF) > (snake.X+rS) || (snake.X-rS) > (food.X+rF) {
		return false
	}
	// we check if they align on the Y axis
	if (food.Y-rF) > (snake.Y+rS) || (snake.Y-rS) > (food.Y+rF) {
		return false
	}

	return true
}