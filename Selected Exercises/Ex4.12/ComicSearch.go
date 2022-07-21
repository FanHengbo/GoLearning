package main

import (
	"fmt"
	"log"

	"Ex4.12/comic"
)

func main() {
	/*
		argc := len(os.Args)
		if argc != 3 {
			log.Fatalln("Usage: ./ComicSearch <get|search> comicNumber / keyowrds")
		}
		newcomic, err := comic.GetComic(2)
		if err != nil {
			log.Fatal(err)
		}
	*/
	//fmt.Println(newcomic.Title, "\t", newcomic.Transcript)

	data := comic.New()
	//data.InitComic()
	err := data.ReadFromFile("data.json")
	if err != nil {
		log.Fatal(err)
	}
	c, _ := data.Get(1)
	fmt.Println(c.Transcript)
	/*
		err := data.SaveToFile("data.json")
		if err != nil {
			log.Fatal(err)
		}
	*/
}
