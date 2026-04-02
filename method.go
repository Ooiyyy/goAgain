package main

import "fmt"

type Hp struct {
	Merk     string
	Tipe     string
	Keluaran int
}

func (h Hp) merk() {
	fmt.Println("Merk hp", h.Merk)
}

func main() {
	hpKu := Hp{
		"xiaomi",
		"12 lite",
		2022,
	}

	hpKu.merk()

}
