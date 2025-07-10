package main

import (
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

	log.Fatal(http.ListenAndServe(appInstance.Options.GetServer(), appInstance.Router))
}
