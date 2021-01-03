package view

import (
	"../model"
	"fmt"
	"html/template"
	"image"
	"image/draw"
	"image/jpeg"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

var (
	playTempl	*template.Template
	VWidth		int = model.Width
	VHeight		int = model.Height
	background	BackgroundImg
)

// Params is the sructure that will be passed to the html template
var Params = struct {
	Width	*int
	Height	*int
}{&VWidth, &VHeight}


//GetHTMLString returns in a string the full html file of the game
func GetHTMLString() string {
	fileInBytes, err := ioutil.ReadFile("view/snake.html") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	fileString := string(fileInBytes) // convert content to a 'string'

	return fileString
}

// init registers the http handlers and creates initial images.
func init() {
	fmt.Println("Init handlers executed")
	snakeHTML := GetHTMLString()
	playTempl = template.Must(template.New("t").Parse(snakeHTML))
	http.HandleFunc("/", loadMainPageHandle)
	http.HandleFunc("/new", newGameHandle)
	http.HandleFunc("/keyPressed", keyPressedHandle)
	http.HandleFunc("/img", imgHandle)
	http.HandleFunc("/score", scoreHandle)

	background = BackgroundImg{}
}

// loadMainPageHandle serves the html page where the user can play.
func loadMainPageHandle(writer http.ResponseWriter, req *http.Request) {
	_ = playTempl.Execute(writer, Params)
}

// newGameHandle signals to start a newgame.
func newGameHandle(writer http.ResponseWriter, req *http.Request) {
	model.ChanNewGame <- true
}

// keyPressedHandle handles the click of a key (in particular an arrow)
func keyPressedHandle(writer http.ResponseWriter, req *http.Request) {
	keyPressed := req.FormValue("code")
	// non-blocking send in case the user stacks up key presses
	select{
	case model.ChanKeyPress <- keyPressed:
		// none
	default:
		fmt.Println("Unable to push into Key Press Channel!")
	}
}

// imgHandle serves images of the player's view
func imgHandle(writer http.ResponseWriter, req *http.Request) {
	img := createGameImg()
	jpeg.Encode(writer, img, nil)

}

// scoreHandle serves images of the player's view
func scoreHandle(writer http.ResponseWriter, req *http.Request) {
	io.WriteString(writer, strconv.Itoa(model.Score)+"&"+strconv.Itoa(model.MaxScore))
}

// createGameImg generates the image of the full game
func createGameImg() image.Image{
	fullImg := image.NewRGBA(background.Bounds())
	draw.Draw(fullImg, background.Bounds(), background, background.Bounds().Min, draw.Src)

	// we first draw the snake
	var snake *model.SnakeNode
	var snakeImg SnakeNodeImg
	snake = model.Snake

	for snake != nil{
		// we draw each of the nodes of the snake i
		snakeImg = SnakeNodeImg{image.Rect(snake.X-model.RadiusSnake,snake.Y-model.RadiusSnake,
				snake.X+model.RadiusSnake,snake.Y+model.RadiusSnake)}

		draw.Draw(fullImg, snakeImg.Bounds(), snakeImg, snakeImg.Bounds().Min, draw.Src)
		snake = snake.Next
	}

	// now we draw the food
	food := model.Food
	foodImg := FoodImg{image.Rect(food.X-model.RadiusFood, food.Y-model.RadiusFood,
		food.X+model.RadiusFood, food.Y+model.RadiusFood)}
	draw.Draw(fullImg, foodImg.Bounds(), foodImg, foodImg.Bounds().Min, draw.Src)

	return fullImg
}