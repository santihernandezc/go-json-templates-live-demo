//go:build js && wasm

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/santihernandezc/go-json/interpreter"
)

func main() {
	c := make(chan struct{}, 0)
	registerCallbacks()
	<-c
}

func registerCallbacks() {
	js.Global().Set("execute", js.FuncOf(execute))
}

func execute(this js.Value, args []js.Value) any {
	if len(args) != 2 {
		return js.ValueOf(newErrorResponse(fmt.Sprintf("Incorrect number of arguments: %d, want 2", len(args))))
	}

	data := make(map[string]any)
	if err := json.Unmarshal([]byte(args[1].String()), &data); err != nil {
		return js.ValueOf(newErrorResponse(fmt.Sprintf("Error parsing data into JSON: %s", err)))
	}

	scanner := interpreter.NewScanner([]byte(args[0].String()))
	tokens := scanner.Scan()
	parser := interpreter.NewParser(tokens)
	statements, err := parser.Parse()
	if err != nil {
		return js.ValueOf(newErrorResponse(fmt.Sprintf("Error parsing template: %s", err)))
	}

	interpreter := interpreter.NewInterpreter(statements, data)
	res, err := interpreter.Interpret()
	if err != nil {
		return js.ValueOf(newErrorResponse(fmt.Sprintf("Error interpreting template: %s", err)))
	}

	return newSuccessResponse(string(res))
}

func newErrorResponse(err string) map[string]any {
	return map[string]any{
		"error": err,
	}
}

func newSuccessResponse(res string) map[string]any {
	return map[string]any{
		"result": res,
	}
}
