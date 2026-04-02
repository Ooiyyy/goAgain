package main

import "fmt"

type User struct {
	Name string
	Age  int
}

type UserService struct{}

func updateuser(u *User) {
	u.Name = "updated"
	u.Age += 5
}
func main() {
	user := User{"budi", 20}

	updateuser(&user)
	fmt.Println(user)
}
