package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Bananas struct {
	amount int
}

type Potatoes struct {
	amount int
}

type Cargo struct {
	BananaAmount int
	PotatoAmount int
}

const MAX_CARGO = 5000

func warehouse(bananaTrucks <-chan Bananas, potatoTrucks <-chan Potatoes,
	cargoTrucks chan Cargo, timer *time.Timer, wg *sync.WaitGroup) {
	pAcc, bAcc := 0, 0
	for {
		select {
		case potatoes := <-potatoTrucks:
			fmt.Printf("\n%v potatoes arrived\n", potatoes.amount)
			if (pAcc + bAcc + potatoes.amount) < MAX_CARGO {
				pAcc += potatoes.amount
			} else {
				fmt.Printf("\nTruck sent\n")
				cargoTrucks <- Cargo{
					BananaAmount: bAcc,
					PotatoAmount: MAX_CARGO - bAcc,
				}
				pAcc, bAcc = (pAcc+bAcc+potatoes.amount)-MAX_CARGO, 0
			}
		case bananas := <-bananaTrucks:
			fmt.Printf("\n%v bananas arrived\n", bananas.amount)
			if (pAcc + bAcc + bananas.amount) < MAX_CARGO {
				bAcc += bananas.amount
			} else {
				fmt.Printf("\nTruck sent\n")
				cargoTrucks <- Cargo{
					BananaAmount: MAX_CARGO - pAcc,
					PotatoAmount: pAcc,
				}
				bAcc, pAcc = (pAcc+bAcc+bananas.amount)-MAX_CARGO, 0
			}
		case truck := <-cargoTrucks:
			fmt.Printf(`
Truck arrived.
Potatoes: %v
Bananas: %v
`, truck.PotatoAmount, truck.BananaAmount)
		case <-timer.C:
			fmt.Printf(`
Warehose remainders:
Potatoes: %v
Bananas: %v
`, pAcc, bAcc)
			wg.Done()
			return
		}
	}
}

func main() {
	timer := time.NewTimer(60 * time.Second)
	var wg sync.WaitGroup
	banaChan, potaChan, carChan := make(chan Bananas), make(chan Potatoes), make(chan Cargo, 1)

	go func() {
		banaTick := time.NewTicker(time.Second)
		for {
			select {
			case <-banaTick.C:
				banaChan <- Bananas{amount: rand.Intn(1000)}
			}
		}
	}()

	go func() {
		potaTick := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-potaTick.C:
				potaChan <- Potatoes{amount: rand.Intn(5000)}
			}
		}
	}()

	wg.Add(1)
	go warehouse(banaChan, potaChan, carChan, timer, &wg)
	wg.Wait()
}
