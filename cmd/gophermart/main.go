package main

import (
	"net/http"

	"github.com/Neeeooshka/gopher-club/internal/wire"
)

func main() {
	appInstance, cleanup, err := wire.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	http.ListenAndServe(appInstance.Options.GetServer(), appInstance.Router)
}
