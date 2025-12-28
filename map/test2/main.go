package main

import (
	"fmt"
)

type item struct {
	name  string
	stock int
}

func main() {
	// Scenario : print inventory report in sorted product name order
	inventory := map[string]item{
		"banana": {name: "Banana", stock: 20},
		"apple":  {name: "apple", stock: 2},
		"orange": {name: "orange", stock: 30},
	}

	keys := make([]string, 0, len(inventory))

	for k := range inventory {
		keys = append(keys, k)
	}

	sortStringsReverse(keys)

	fmt.Println("Inventory report (sorted bu product code):")
	for _, k := range keys {
		item := inventory[k]
		fmt.Printf("- (%s) : %d \n", item.name, item.stock)
	}
}
func sortStrings(keys []string) {
	for i := 1; i < len(keys); i++ {
		current := keys[i]
		j := i - 1
		for j >= 0 && keys[j] > current {
			keys[j+1] = keys[j]
			j--
		}
		keys[j+1] = current
	}
}

func sortStringsReverse(keys []string) {
	for i := 1; i < len(keys); i++ {
		current := keys[i]
		j := i - 1
		for j >= 0 && keys[j] < current {
			keys[j+1] = keys[j]
			j--
		}
		keys[j+1] = current

	}

}
