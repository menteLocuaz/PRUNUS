package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "Admin123"
	hash := "$2a$10$U.sUS/qwAXlDPrJZ9wAaLe78DmRtcnWVY39wFp85YLiL0iIVPVkkK"
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Password matches hash")
	}
}
