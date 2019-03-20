package main

import "testing"
import "fmt"
import "strings"
import "errors"

//import "net/http/httptest"

type hashTestCase struct {
	password string
	hash     string
}

func TestHashPassword(t *testing.T) {
	passwordHashes := []hashTestCase{
		{"password", "5BAA61E4C9B93F3F0682250B6CF8331B7EE68FD8"},
		{"test", "A94A8FE5CCB19BA61C4C0873D391E987982FBBD3"},
	}
	for _, testParameters := range passwordHashes {
		t.Run(fmt.Sprintf("%s", testParameters.password), func(t *testing.T) {
			singlePasswordTest(t, testParameters.password, testParameters.hash)
		})
	}
}

func singlePasswordTest(t *testing.T, password string, expectedHash string) {
	hash := hashPassword(password)
	if expectedHash != hash {
		t.Errorf("Hash of password was expected to be %s, was %s", expectedHash, hash)
	}
}

func testBuildUrl(t *testing.T) {
	fullHash := "5BAA61E4C9B93F3F0682250B6CF8331B7EE68FD8"
	expectedUrl := "https://api.pwnedpasswords.com/range/5BAA6"
	url := buildUrl(fullHash)
	if url != expectedUrl {
		t.Fatalf("url should have been %s, was %s", expectedUrl, url)
	}
}

type TestClient struct {
	results    []string
	shouldFail bool
}

func (t TestClient) Get(hash string) (response string, err error) {
	if t.shouldFail {
		return "", errors.New("client failed")
	}
	return strings.Join(t.results, "\n"), nil
}

func TestQueryApi(t *testing.T) {
	results := []string{
		"00001111112222223333333444444555556:8",
		"1E4C9B93F3F0682250B6CF8331B7EE68FD8:8",
	}
	client := &TestClient{results, false}
	count, err := queryApi(client, "5BAA61E4C9B93F3F0682250B6CF8331B7EE68FD8")
	if err != nil {
		t.Fatalf("An error occured: %v", err)
	}
	if count != 8 {
		t.Fatalf("Hash should be found 8 times, was found %v times", count)
	}
}

func TestQueryApiNoResult(t *testing.T) {
	results := []string{
		"00001111112222223333333444444555556:8",
		"1E4C9B93F3F0682250B6CF8331B7EE68FD8:8",
	}
	client := &TestClient{results, false}
	count, err := queryApi(client, "0000000000000000000000000000000000000000")
	if err != nil {
		t.Fatalf("An error occured: %v", err)
	}
	if count != 0 {
		t.Fatalf("Hash should not be found, was found %v times", count)
	}
}

func TestQueryApiNetworkFail(t *testing.T) {
	results := []string{}
	client := &TestClient{results, true}
	_, err := queryApi(client, "5BAA61E4C9B93F3F0682250B6CF8331B7EE68FD8")
	if err == nil {
		t.Fatal("An error should have occured")
	}
}
