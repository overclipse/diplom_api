package main

import (
	"diplom_api/internal/api_reqest"
	"sync"
)

func main() {
	wg := new(sync.WaitGroup)
	wg.Go(func() {
		api_reqest.Fsssp()
	})
	wg.Wait()
}