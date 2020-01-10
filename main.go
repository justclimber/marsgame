package main

import (
	"aakimov/marslang/lexer"
	"aakimov/marslang/object"
	"aakimov/marslang/parser"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func setupRoutes() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/ws", wsEndpoint)
	http.HandleFunc("/save_source_code", saveSourceCodeEndpoint)
}

func saveSourceCodeEndpoint(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	sourceCode := bytes.Trim(reqBody, "\"")
	sourceCodeStr := string(sourceCode)
	sourceCodeStr = strings.ReplaceAll(sourceCodeStr, "\\n", "\n")
	fmt.Printf("Program to parse: %s\n", sourceCodeStr)

	l := lexer.New(sourceCodeStr)
	p := parser.New(l)
	astProgram, err := p.Parse()
	if err != nil {
		fmt.Printf("Parsing error: %s\n", err.Error())
	}
	env := object.NewEnvironment()
	_, err = astProgram.Exec(env)
	if err != nil {
		fmt.Printf("Runtime error: %s\n", err.Error())
	}

	vars, err := env.GetVarsAsJson()
	if err != nil {
		fmt.Printf("Marshaling error: %s\n", err.Error())
	}

	w.Write(vars)

	env.Print()
	fmt.Printf("Returned to client: %s", string(vars))

}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")

	reader(ws)
}

func main() {
	fmt.Println("Running a http server")
	setupRoutes()
	http.ListenAndServe(":80", nil)
}
