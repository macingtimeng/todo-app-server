package main

import (
	"runtime"
	"todo-app/handler"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	handler.StartApp()
}
