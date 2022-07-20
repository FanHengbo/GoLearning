package comic

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Comic struct {
	Month      string
	Num        int
	Link       string
	Year       string
	Transcript string
	Title      string
}

var sema = make(chan struct{}, 30)

type ComicsDatabase struct {
	database Comics
}

type ComicNum int
type Comics map[ComicNum]*Comic

func New() *ComicsDatabase {
	return &ComicsDatabase{database: make(Comics)}
}

var Database *ComicsDatabase = New()

func (d *ComicsDatabase) Get(num ComicNum) (*Comic, error) {
	item, _ := d.database[num]
	return item, nil
}
func init() {
	fmt.Println("Database initializing...")
	var n sync.WaitGroup
	for i := 0; i < 100; i++ {
		n.Add(1)
		go func() {
			newComic, err := GetComic(ComicNum(i))
			if err != nil {
				log.Print(err)
			}
			Database.database[ComicNum(newComic.Num)] = newComic
			n.Done()
		}()
	}
	n.Wait()
	fmt.Println("Database Done")
}
func GetComicQuantity() (int, error) {
	url := "https://xkcd.com/info.0.json"
	resp, err := http.Get(url)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return -1, err
	}
	result := struct {
		Num int
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return -1, err
	}
	return result.Num, nil

}
func GetComic(num ComicNum) (*Comic, error) {
	url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", num)
	sema <- struct{}{}
	resp, err := http.Get(url)
	<-sema
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	var result Comic
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
