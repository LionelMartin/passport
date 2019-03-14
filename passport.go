package main

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("The program should be given a password as argument")
		os.Exit(1)
	}
	password := os.Args[1]
	hash := hashPassword(password)
	count, err := queryApi(hash)
	if err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		fmt.Printf("The password was found %d times\n", count)
	} else {
		fmt.Println("The password was not found")
	}
}

func hashPassword(password string) string {
	data := []byte(password)
	return strings.ToUpper(fmt.Sprintf("%x", sha1.Sum(data)))
}

func queryApi(hash string) (int, error) {
	hashSizeSent := 5
	url := "https://api.pwnedpasswords.com/range/"
	url = strings.Join([]string{url, string(hash[0:hashSizeSent])}, "")
	fmt.Printf("Calling: %v\n", url)
	res, err := http.Get(url)
	if err != nil {
		return -1, err
	}
	defer res.Body.Close()
	scanner := bufio.NewScanner(res.Body)
	if err != nil {
		return -1, err
	}
	count := 0
	for count == 0 && scanner.Scan() {
		text := strings.Split(scanner.Text(), ":")
		if text[0] == string(hash[hashSizeSent:]) {
			count, err = strconv.Atoi(text[1])
			if err != nil {
				return -1, err
			}
		}
	}
	return count, nil
}
