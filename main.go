package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
)

type Latest struct {
	Assets []struct {
		Name               string `json:"name"`
		ContentType        string `json:"content_type"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func main() {
	flag.Parse()
	args := flag.Args()

	var repo string
	if len(args) == 0 {

		fmt.Println("input onwer/repository:")
		var sc = bufio.NewScanner(os.Stdin)
		sc.Scan()
		repo = sc.Text()
	} else {
		repo = args[0]
	}

	var latest Latest
	json.Unmarshal(getbytes(repo), &latest)

	for _, asset := range latest.Assets {

		if strings.Contains(strings.ToLower(asset.BrowserDownloadURL), runtime.GOOS) {

			fmt.Printf("found %s\n", asset.BrowserDownloadURL)

			resp, err := http.Get(asset.BrowserDownloadURL)

			if err != nil {
			}
			defer resp.Body.Close()

			var path string
			if len(args) == 2 {
				path = args[1] + "/" + asset.Name
			} else {
				path = "./" + asset.Name
			}

			file, err := os.Create(path)
			if err != nil {
			}
			defer file.Close()

			byteArray, _ := ioutil.ReadAll(resp.Body)
			file.Write(byteArray)

			fmt.Printf("dowload done : %s\n", path)

			break
		} else {
			fmt.Printf("miss match %s\n", asset.BrowserDownloadURL)
		}
	}
}

func getbytes(owerrepo string) []byte {

	url := "https://api.github.com/repos/" + owerrepo + "/releases/latest"
	fmt.Println(url)
	resp, err := http.Get(url)

	if err != nil {
		defer resp.Body.Close()
	}

	byteArray, _ := ioutil.ReadAll(resp.Body)

	return byteArray
}
