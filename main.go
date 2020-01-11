package main

import (
	"aakimov/marslang/interpereter"
	"aakimov/marslang/lexer"
	"aakimov/marslang/object"
	"aakimov/marslang/parser"
	"bytes"
	"encoding/json"
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
	p, err := parser.New(l)
	if err != nil {
		respondWithError(err.Error(), "Lexing", w)
		return
	}
	astProgram, err := p.Parse()
	if err != nil {
		respondWithError(err.Error(), "Parsing", w)
		return
	}
	env := object.NewEnvironment()
	_, err = interpereter.Exec(astProgram, env)
	if err != nil {
		respondWithError(err.Error(), "Runtime", w)
		return
	}

	vars, err := env.GetVarsAsJson()
	if err != nil {
		respondWithError(err.Error(), "Marshaling", w)
		return
	}

	_, err = w.Write(vars)
	if err != nil {
		respondWithError(err.Error(), "Responding", w)
		return
	}

	env.Print()
	fmt.Printf("Returned to client: %s", string(vars))
}

func respondWithError(msg, prefix string, w http.ResponseWriter) {
	fmt.Printf("%s error: %s\n", prefix, msg)
	_, err := w.Write(errorToJson(msg))
	if err != nil {
		fmt.Printf("Writing to response error: %s\n", err.Error())
	}
}

func errorToJson(msg string) []byte {
	errJson := make(map[string]string)
	errJson["error"] = msg
	errBytes, _ := json.Marshal(errJson)

	return errBytes
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
