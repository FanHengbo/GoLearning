package comic

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
type Comics map[ComicNum]Comic

func New() *ComicsDatabase {
	return &ComicsDatabase{database: make(Comics)}
}

func (d *ComicsDatabase) Get(num ComicNum) (Comic, error) {
	item := d.database[num]
	return item, nil
}
func (d *ComicsDatabase) SaveToFile(fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	fmt.Println("Saving data to file...")
	encoder := json.NewEncoder(file)
	err = encoder.Encode(d.database)
	if err != nil {
		return fmt.Errorf("writing error")
	}
	fmt.Println("Data saved")
	return nil
}
func (d *ComicsDatabase) ReadFromFile(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}
func (d *ComicsDatabase) InitComic() {
	type Item struct {
		comicInfo Comic
		num       ComicNum
		err       error
	}
	//Replace 50 with maximum comics count later
	itemChannel := make(chan Item)
	fmt.Println("Database initializing...")
	var i ComicNum
	var n sync.WaitGroup
	// Not passing i into goroutine is so stupid...
	for i = 0; i < 50; i++ {
		n.Add(1)
		go func(num ComicNum) {
			var it Item
			it.comicInfo, it.err = GetComic(num)
			it.num = num
			itemChannel <- it
			n.Done()
		}(i)
	}

	go func() {
		n.Wait()
		close(itemChannel)
	}()
	for it := range itemChannel {
		if it.err != nil {
			log.Fatal("http.get error")
		}
		d.database[it.num] = it.comicInfo
	}
	fmt.Println("Database has already initialized")
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
func GetComic(num ComicNum) (Comic, error) {
	var result Comic
	url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", num)
	sema <- struct{}{}
	resp, err := http.Get(url)
	<-sema
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return result, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, err
	}
	return result, nil
}
