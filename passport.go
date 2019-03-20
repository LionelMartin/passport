package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var apiHashSizeSent = 5
var apiUrl = "https://api.pwnedpasswords.com/range/%s"

func main() {
	if len(os.Args) != 2 {
		fmt.Println("The program should be given a password as argument")
		os.Exit(1)
	}
	password := os.Args[1]
	hash := hashPassword(password)

	var apiClient = &RealApiClient{}
	count, err := queryApi(apiClient, hash)
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

type ApiClient interface {
	Get(hash string) (string, error)
}

type RealApiClient struct {
}

func (client *RealApiClient) Get(hash string) (string, error) {
	url := buildUrl(hash)
	fmt.Printf("Calling: %v\n", url)
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func buildUrl(hash string) string {
	return fmt.Sprintf(apiUrl, string(hash[0:apiHashSizeSent]))
}

func queryApi(client ApiClient, hash string) (int, error) {
	response, err := client.Get(hash)
	if err != nil {
		return 0, err
	}
	for _, line := range strings.Split(response, "\n") {
		text := strings.Split(line, ":")
		if text[0] == string(hash[apiHashSizeSent:]) {
			count, err := strconv.Atoi(text[1])
			if err != nil {
				return 0, err
			}
			return count, nil
		}
	}
	return 0, nil
}
