package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

// InventoryCache is a real-world style in-memory stock cache used by an API.
// It is safe for concurrent read/write access.
type InventoryCache struct {
	mu    sync.RWMutex
	stock map[string]int
}

func NewInventoryCache() *InventoryCache {
	return &InventoryCache{
		stock: make(map[string]int),
	}
}

func (c *InventoryCache) GetSKU(sku string) (int, bool) {
	c.mu.RLock()
	qty, ok := c.stock[sku]
	c.mu.RUnlock()
	return qty, ok
}

func (c *InventoryCache) SetSKU(sku string, qty int) {
	c.mu.Lock()
	c.stock[sku] = qty
	c.mu.Unlock()
}

func (c *InventoryCache) AdjustSKU(sku string, delta int) {
	c.mu.Lock()
	c.stock[sku] += delta
	c.mu.Unlock()
}

// UnsafeInventory is intentionally unsafe and can panic under concurrent access.
type UnsafeInventory struct {
	stock map[string]int
}

func NewUnsafeInventory() *UnsafeInventory {
	return &UnsafeInventory{stock: make(map[string]int)}
}

func (u *UnsafeInventory) GetSKU(sku string) (int, bool) {
	qty, ok := u.stock[sku]
	return qty, ok
}

func (u *UnsafeInventory) SetSKU(sku string, qty int) {
	u.stock[sku] = qty
}

func (u *UnsafeInventory) AdjustSKU(sku string, delta int) {
	u.stock[sku] += delta
}

func runSafeDemo() {
	cache := NewInventoryCache()
	cache.SetSKU("sku-1001", 50)

	var wg sync.WaitGroup

	// Simulate order workers decrementing stock concurrently.
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				cache.AdjustSKU("sku-1001", -1)
				time.Sleep(2 * time.Millisecond)
			}
		}()
	}

	// Simulate API reads for stock display.
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				if qty, ok := cache.GetSKU("sku-1001"); ok {
					_ = qty
				}
				time.Sleep(2 * time.Millisecond)
			}
		}()
	}

	wg.Wait()

	qty, _ := cache.GetSKU("sku-1001")
	fmt.Println("safe demo final stock:", qty)
}

func runUnsafeDemo() {
	cache := NewUnsafeInventory()
	cache.SetSKU("sku-1001", 50)

	var wg sync.WaitGroup

	// This is a realistic workload but unsafe for a plain map.
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				cache.AdjustSKU("sku-1001", -1)
				time.Sleep(2 * time.Millisecond)
			}
		}()
	}

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				if qty, ok := cache.GetSKU("sku-1001"); ok {
					_ = qty
				}
				time.Sleep(2 * time.Millisecond)
			}
		}()
	}

	wg.Wait()

	qty, _ := cache.GetSKU("sku-1001")
	fmt.Println("unsafe demo final stock:", qty)
}

func main() {
	mode := flag.String("mode", "safe", "demo mode: safe or unsafe")
	flag.Parse()

	switch *mode {
	case "safe":
		runSafeDemo()
	case "unsafe":
		runSafeDemo()
	default:
		fmt.Println("unknown mode:", *mode)
	}
}
