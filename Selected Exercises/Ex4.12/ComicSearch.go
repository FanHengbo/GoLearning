package main

import (
	"fmt"

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
	fmt.Println(comic.Database.Get(0))
}
