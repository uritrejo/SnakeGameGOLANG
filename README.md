# SnakeGameGOLANG
This is a basic implementation of the game Snake. It is played from a browser. All logic and image generation is written in GoLang.  A bit of HTML, CSS and Javascript were also used to render the view.

## Usage
The game can be run by cloning this repo and executing "go run main.go" in the main directory. This will create the server and prompt a new tab on your default browser.

## Design
Follows principles of Model-View-Controller. 
Hosted by a webserver. Javascript code refreshes the board image and scores at a specified frame rate. The controller runs on a separate go-routine from the main server.
