/* (c) 2021 Davis Tibbz. MIT License https://github.com/longwater1234       */
package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// global error handler
func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

var (
	FOUND    uint32 = 0
	NOTFOUND uint32 = 0
)

// console colors
const (
	RED   string = "\033[31m"
	GREEN string = "\033[32m"
	RESET string = "\033[0m"
)

type Website struct {
	Url string `json:"url"`
}

var WebsiteArr []Website

// Search :for the given username in the website
func Search(wg *sync.WaitGroup, c *http.Client, w Website, username string) {
	var finalUrl string = strings.ReplaceAll(w.Url, "%", username)
	mama := strings.SplitAfter(w.Url, "//")[1]
	defer wg.Done()
	res, err := c.Get(finalUrl)
	if err != nil {
		fmt.Printf("[!] failed on %s \n", mama)
		return
	}
	defer res.Body.Close()

	var exists string
	if res.StatusCode == 200 {
		exists = string(GREEN) + "\u2713"
		atomic.AddUint32(&FOUND, 1)
	} else {
		exists = string(RED) + "x"
		atomic.AddUint32(&NOTFOUND, 1)
	}

	fmt.Printf("%v %s on %s? %v \n", exists, username, mama, string(RESET))
}

func main() {
	var wg sync.WaitGroup

	var username string = os.Args[1]
	if len(username) < 2 {
		check(errors.New("username is too short"))
	}

	reg := regexp.MustCompile("^[a-zA-Z0-9_-]{2,}$")

	if !reg.MatchString(username) {
		check(errors.New("username is invalid"))
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	starting := time.Now()
	fmt.Println("Starting search...")

	//open the json file
	f, err := os.Open("./websites.json")
	check(err)
	defer f.Close()

	//decode the JSON file to a Slice
	r := bufio.NewReader(f)
	jd := json.NewDecoder(r)
	err1 := jd.Decode(&WebsiteArr)
	check(err1)

	//do the Search for each site
	for _, w := range WebsiteArr {
		wg.Add(1)
		go Search(&wg, &client, w, username)
	}

	fmt.Println("Main: Waiting for workers to finish")
	wg.Wait()
	fmt.Printf("Search completed in: %d ms \n", time.Since(starting).Milliseconds())
	fmt.Print("\n")
	fmt.Printf("%s found in %d SITES \n", username, FOUND)
	fmt.Printf("%s NOT found in %d SITES \n", username, NOTFOUND)

}
