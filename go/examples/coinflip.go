package main

import (
	. "chatterbet/coinflip"
	"fmt"
	"sync"
	"time"
)

type Context struct {
	Game *Coinflip

	gameMessageOutput <-chan string

	done           chan bool
	gameToStdout   chan string
	playerToGame   chan Message
	playerToStdout chan Message

	playersDone bool
	gameDone    bool

	WaitGroup sync.WaitGroup
}

func NewContext() *Context {
	ctx := &Context{}

	ctx.gameToStdout = make(chan string, 10)
	ctx.playerToGame = make(chan Message, 10)
	ctx.playerToStdout = make(chan Message, 10)

	ctx.Game = NewCoinflip()
	ctx.Game.MessageChannel = ctx.gameToStdout

	return ctx
}

func (c *Context) SendMessage(m Message) {
	c.playerToGame <- m
	c.playerToStdout <- m
}

func ExecuteGame(ctx *Context) {

	ctx.WaitGroup.Done()
	ctx.gameDone = true
}

func ExecuteInteractions(ctx *Context) {
	ctx.SendMessage(Message{
		User:    "bob",
		Content: "!play coinflip",
	})
	time.Sleep(1 * time.Second)
	ctx.SendMessage(Message{
		User:    "bob",
		Content: "!bet head",
	})
	time.Sleep(2 * time.Second)
	ctx.SendMessage(Message{
		User:    "alice",
		Content: "!bet tail",
	})

	ctx.WaitGroup.Done()
	ctx.playersDone = true
}

func ExecuteStdout(ctx *Context) {
	for {
		done := false
		select {
		case msg := <-ctx.playerToStdout:
			fmt.Printf("(stdout) %s: %s\n", msg.User, msg.Content)
		case str := <-ctx.gameToStdout:
			fmt.Printf("(stdout) %s: %s\n", "chatterbet", str)
		default:
			if ctx.gameDone && ctx.playersDone {
				done = true
			}
		}
		if done {
			break
		}
	}

	ctx.WaitGroup.Done()
}

func main() {
	fmt.Println("\n(main) Running coinflip example...")

	ctx := NewContext()

	ctx.WaitGroup.Add(3)

	go ExecuteInteractions(ctx)
	go ExecuteGame(ctx)
	go ExecuteStdout(ctx)

	ctx.WaitGroup.Wait()

	fmt.Println("(main) Done.")
}
