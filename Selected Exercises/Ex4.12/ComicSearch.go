package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"Ex4.12/comic"
)

func main() {
	searchNum, _ := strconv.Atoi(os.Args[2])
	data := comic.New()
	err := data.ReadFromFile("data.json")
	if err != nil {
		log.Fatal(err)
	}
	c, _ := data.Get(comic.ComicNum(searchNum))
	fmt.Println(c.Transcript)
}
