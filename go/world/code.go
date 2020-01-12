package world

import (
	"aakimov/marslang/ast"
	"aakimov/marslang/interpereter"
	"aakimov/marslang/lexer"
	"aakimov/marslang/object"
	"aakimov/marslang/parser"

	"encoding/json"
	"log"
	"sync"
	"time"
)

type Code struct {
	id         string
	env        *object.Environment
	mu         sync.Mutex
	astProgram *ast.StatementsBlock
	outputCh   chan *MechOutputVars
}

type MechOutputVars struct {
	MThrottle float64
	RThrottle float64
}

func newMechOutputVarsFromEnv(env *object.Environment) *MechOutputVars {
	return &MechOutputVars{
		MThrottle: getFloatVarFromEnv("mThrottle", env),
		RThrottle: getFloatVarFromEnv("rThrottle", env),
	}
}

func getFloatVarFromEnv(varName string, env *object.Environment) float64 {
	envObj, ok := env.Get(varName)
	if !ok {
		return 0
	}
	objFloat, ok := envObj.(*object.Float)
	if !ok {
		log.Fatalf("%s should be type float, %s given", varName, objFloat.Type())
	}
	return objFloat.Value
}

func (c *Code) Run() {
	ticker := time.NewTicker(2 * time.Second)

	// endless loop here
	for _ = range ticker.C {
		log.Printf("Code run tick\n")
		c.listenTheWorld()
		if c.astProgram == nil {
			// waiting for the fist program
			continue
		}
		env := object.NewEnvironment()

		c.mu.Lock()
		astProgram := c.astProgram
		c.mu.Unlock()

		_, err := interpereter.Exec(astProgram, env)
		if err != nil {
			log.Printf("Runtime error: %s", err.Error())
		}

		c.outputCh <- newMechOutputVarsFromEnv(env)

		env.Print()
	}
}

func (c *Code) listenTheWorld() {
	// select for channels
}

func (c *Code) saveAst(ast *ast.StatementsBlock) {
	c.mu.Lock()
	c.astProgram = ast
	c.mu.Unlock()
}

func ParseSourceCode(sourceCode string) (*ast.StatementsBlock, []byte) {
	log.Printf("Program to parse: %s\n", sourceCode)

	l := lexer.New(sourceCode)
	p, err := parser.New(l)
	if err != nil {
		return nil, respondWithError(err.Error(), "Lexing")
	}
	astProgram, err := p.Parse()
	if err != nil {
		return nil, respondWithError(err.Error(), "Parsing")
	}

	return astProgram, nil
}

func respondWithError(msg, prefix string) []byte {
	log.Printf("%s error: %s\n", prefix, msg)
	return errorToJson(msg)
}

func errorToJson(msg string) []byte {
	errJson := make(map[string]string)
	errJson["error"] = msg
	errBytes, _ := json.Marshal(errJson)

	return errBytes
}
