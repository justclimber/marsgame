package world

import (
	"github.com/justclimber/marsgame/server"
	"github.com/justclimber/marsgame/wal"
	"github.com/justclimber/marslang/ast"
	"github.com/justclimber/marslang/lexer"
	"github.com/justclimber/marslang/parser"

	"log"
	"time"
)

type Player struct {
	id             uint32
	world          *World
	client         *server.Client
	mech           *Mech
	mainProgram    *Code
	runSpeedMs     time.Duration
	codeCommandsCh chan *Commands
	io4ClientCh    chan *IO4Client
	codeSaveCh     chan *ast.StatementsBlock
	flowCh         chan ProgramState
	errorCh        chan *Error
	commandsCh     chan *server.Command
	terminateCh    chan bool
	wal            *wal.ObjectObserver
}

func NewPlayer(
	id uint32,
	client *server.Client,
	w *World, m *Mech,
	runSpeedMs time.Duration,
	objectLogManager *wal.ObjectObserver,
) *Player {
	player := &Player{
		id:             id,
		world:          w,
		client:         client,
		mech:           m,
		mainProgram:    NewCode("main"),
		runSpeedMs:     runSpeedMs,
		codeCommandsCh: make(chan *Commands, 3),
		codeSaveCh:     make(chan *ast.StatementsBlock, 1),
		io4ClientCh:    make(chan *IO4Client, 1),
		flowCh:         make(chan ProgramState, 1),
		errorCh:        make(chan *Error, 10),
		commandsCh:     make(chan *server.Command, 10),
		terminateCh:    make(chan bool, 1),
		wal:            objectLogManager,
	}
	return player
}

func (p *Player) saveAstCode(sourceCode string) {
	l := lexer.New(sourceCode)
	pr, err := parser.New(l)
	if err != nil {
		p.errorCh <- &Error{
			ErrorType: Lexing,
			Message:   err.Error(),
		}
		log.Printf("Lexing error: %s", err.Error())
		return
	}
	astProgram, err := pr.Parse()
	if err != nil {
		p.errorCh <- &Error{
			ErrorType: Parsing,
			Message:   err.Error(),
		}
		log.Printf("Parsing error: %s", err.Error())
		return
	}
	log.Println("Code parsed")
	p.codeSaveCh <- astProgram
}

func (p *Player) listen() {
	log.Printf("Player [%d] listening started", p.id)
	for {
		select {
		case commands := <-p.codeCommandsCh:
			p.mech.Lock()
			p.mech.commands = commands
			p.mech.Unlock()
		case codeError := <-p.errorCh:
			p.client.PackAndSendCommand("codeError", codeError)
		case io4Client := <-p.io4ClientCh:
			p.client.PackAndSendCommand("codeInputOutput", io4Client)
		case <-p.terminateCh:
			return
		case command := <-p.commandsCh:
			switch command.Type {

			}
		}
	}
}
