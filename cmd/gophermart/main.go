package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Neeeooshka/gopher-club/internal/wire"
)

func main() {
	appInstance, cleanup, err := wire.InitializeApp()
	if err != nil {
		log.Fatalf("не удалось инициализировать приложение: %v", err)
	}
	defer cleanup()

	err = http.ListenAndServe(appInstance.Options.GetServer(), appInstance.Router)
	if err != nil {
		panic(fmt.Errorf("error starting server: %s", err))
	}
}
