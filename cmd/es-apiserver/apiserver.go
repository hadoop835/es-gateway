package main

import "github.com/es-gateway/cmd/es-apiserver/app"

func main() {
	cmd := app.WithApiServerCommand()
	err := cmd.Execute()
	if err != nil {

	}
}
