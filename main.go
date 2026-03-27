package main

import "fmt"

func main() {
	nama := "Budi"
	umur := 10

	fmt.Println(nama)
	fmt.Println(umur)
	LPersegi(4)

	user := User{
		Nama: "agus",
		Umur: 17,
	}
	fmt.Println(user.Nama)
	fmt.Println(user.Umur)

	// ubahUmur(&user, 10)
	tambahUmur(user, 5)
	tambahUmurPointer(&user, 5)
}

func LPersegi(sisi int) {
	L := sisi * sisi
	fmt.Println(L)
}

// type User struct {
// 	Nama string
// 	Umur int
// }

func ubahUmur(u *User, umurBaru int) {
	u.Umur = umurBaru
	fmt.Println("Umur Agus ternyata", umurBaru)
}

func tambahUmur(u User, tambahanUmur int) {
	u.Umur = u.Umur + tambahanUmur
	fmt.Println("umurnya agus jadi", u.Umur)
}

func tambahUmurPointer(u *User, tambahanUmur int) {
	u.Umur = u.Umur + tambahanUmur
	fmt.Println("umurnya agus jadi", u.Umur)
}
