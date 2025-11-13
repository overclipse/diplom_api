package main

import (
	"diplom_api/internal/apireqest"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	 finalURL := "https://api.damia.ru/fssp/isps?inn=7712040126&format=1&key=2268a80e1f11a48f8657c69c71f1c00d41be9219"
	 rosstatURL := "https://api.damia.ru/rs/balance?inn=7712040126&key=67266ba78d7779083310826cc491faf420858d1f"
	wg.Go(func() {
		apireqest.Fsssp(finalURL)
	})
	wg.Go(func() {
		apireqest.Rosstat(rosstatURL)
	})
	wg.Wait()
}
