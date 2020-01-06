# Get Reddit video from a given url

Give it a url, it will try to scrap the video and send you back.

## Usage

Chat to the bot, then

	/reddit reddit_link

## Installation

1. Clone this repo.

	You may want to install dependencies:
	```
	go get gopkg.in/tucnak/telebot.v2
	go get github.com/buger/jsonparser
	```
	
	Since this repo is 'pre-config' as a module for another go app, if you want to make it stand-alone instance, change the main package. 
	
	```
	redditvid.go
	
	package main
	         ^
		 |
		 This one
	```
	
	And a workaround to hang the app: add `select {}` in the end of the main fn.
	```
	func main() {
		go RedditVideoBot()
		select {} 
	}    
	```
	
2. Run `go install`.
3. `export SECRET_TOKEN=your_telegram_bot_token`.
4. Run `$GOPATH/bin/redditvid`.
