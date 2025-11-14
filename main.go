package main

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"diplom_api/internal/apireqest"
	"diplom_api/internal/ui"
)

func main() {
	cfg, err := ui.ConfigureRequests()
	if err != nil {
		if errors.Is(err, ui.ErrMenuAborted) {
			fmt.Println("Запросы отменены.")
			return
		}
		log.Fatalf("не удалось запустить меню настроек: %v", err)
	}

	urls, err := cfg.URLs()
	if err != nil {
		log.Fatalf("не удалось собрать URL: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		apireqest.Fsssp(urls.Fssp)
	}()

	go func() {
		defer wg.Done()
		apireqest.Rosstat(urls.Rosstat)
	}()

	go func() {
		defer wg.Done()
		apireqest.Arb(urls.Arb)
	}()

	wg.Wait()
}
