package main

import "fmt"

type Item struct {
	Name  string
	Stock int
}

func main() {
	// Scenario: print inventory report in sorted product name order.
	inventory := map[string]Item{
		"banana": {Name: "Banana", Stock: 12},
		"apple":  {Name: "Apple", Stock: 18},
		"orange": {Name: "Orange", Stock: 7},
	}

	keys := make([]string, 0, len(inventory))
	for k := range inventory {
		keys = append(keys, k)
	}
	sortStrings(keys)

	fmt.Println("Inventory report (sorted by product code):")
	for _, k := range keys {
		item := inventory[k]
		fmt.Printf("- %s (%s): %d\n", k, item.Name, item.Stock)
	}
}

// sortStrings sorts a string slice in ascending order using insertion sort.
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

// sortStringsDesc sorts a string slice in descending order using insertion sort.
func sortStringsDesc(keys []string) {
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
