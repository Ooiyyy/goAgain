package main

import (
	"fmt"
)

type Product struct {
	Name  string
	Price int
}

type ProductService struct{}

func (p Product) GetInfo() (string, int) {
	return p.Name, p.Price
}

func (p *Product) Discount(percent int) {
	diskon := p.Price * percent / 100
	p.Price = p.Price - diskon
}

func (p Product) IsExpensive() bool {
	if p.Price > 5000 {
		return true
	}
	return false
}

func main() {
	produk1 := Product{"buku", 4000}
	produk1.Discount(10)
	fmt.Println(produk1.Price)
	fmt.Println(produk1.IsExpensive())
}
