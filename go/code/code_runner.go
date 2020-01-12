package code

import (
	"aakimov/marslang/interpereter"
	"aakimov/marslang/lexer"
	"aakimov/marslang/object"
	"aakimov/marslang/parser"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func SaveSourceCode(request []byte, w http.ResponseWriter) {
	sourceCode := bytes.Trim(request, "\"")
	sourceCodeStr := string(sourceCode)
	sourceCodeStr = strings.ReplaceAll(sourceCodeStr, "\\n", "\n")
	log.Printf("Program to parse: %s\n", sourceCodeStr)

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
	log.Printf("Returned to client: %s\n", string(vars))
}

func respondWithError(msg, prefix string, w http.ResponseWriter) {
	log.Printf("%s error: %s\n", prefix, msg)
	_, err := w.Write(errorToJson(msg))
	if err != nil {
		log.Printf("Writing to response error: %s\n", err.Error())
	}
}

func errorToJson(msg string) []byte {
	errJson := make(map[string]string)
	errJson["error"] = msg
	errBytes, _ := json.Marshal(errJson)

	return errBytes
}
