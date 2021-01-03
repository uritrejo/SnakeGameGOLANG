package view

import (
	"image"
	"image/color"
	"../model"
)

var (
	pink color.Color = color.RGBA{255,204,255, 1}
	green color.Color = color.RGBA{0, 204, 102, 1}
	gray color.Color = color.RGBA{96, 96, 96, 1}
	redDead color.Color = color.RGBA{132, 51, 51, 1}
)

// BackgroundImg will implement the image interface to create the game's view
type BackgroundImg struct {}

func (b BackgroundImg) ColorModel() color.Model {
	return color.RGBAModel
}

func (b BackgroundImg) Bounds() image.Rectangle {
	return image.Rect(0,0, model.Width, model.Height)
}

func (b BackgroundImg) At(x, y int) color.Color {
	return pink
}

// SnakeNodeImg will contain the image of a node of the snake
type SnakeNodeImg struct {
	Bounds_	image.Rectangle
}

func (s SnakeNodeImg) ColorModel() color.Model {
	return color.RGBAModel
}

func (s SnakeNodeImg) Bounds() image.Rectangle {
	return s.Bounds_
}

func (s SnakeNodeImg) At(x, y int) color.Color {
	if model.Dead {
		return redDead
	}else {
		return green
	}
}

// FoodImg will contain the image of a node of the snake
type FoodImg struct {
	Bounds_	image.Rectangle
}

func (s FoodImg) ColorModel() color.Model {
	return color.RGBAModel
}

func (s FoodImg) Bounds() image.Rectangle {
	return s.Bounds_
}

func (s FoodImg) At(x, y int) color.Color {
	return gray
}