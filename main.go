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

	if cfg.RunFssp {
		wg.Go(func() {
			apireqest.Fsssp(urls.Fssp)
		})
	}

	if cfg.RunRosstat {
		wg.Go(func() {
			apireqest.Rosstat(urls.Rosstat)
		})
	}

	if cfg.RunArb {
		wg.Go(func() {
			apireqest.Arb(urls.Arb)
		})
	}

	wg.Wait()
}
