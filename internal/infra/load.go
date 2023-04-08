package infra

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

type JSCooldown struct {
	Threads int
	Replies int
	Images  int
}

type JSBoard struct {
	Board             string // code
	Title             string // name
	Ws_board          int
	Per_page          int
	Pages             int
	Max_filesize      int
	Max_comment_chars int
	Image_limit       int
	Cooldowns         JSCooldown
	Meta_description  string // description
	Is_archived       int
}

type JSThread struct {
	No            int
	Last_modified int
	Replies       int
}

type JSPost struct {
	No           int // ID
	Resto        int // 0 if post is thread
	Sticky       int
	Closed       int
	Now          string // Created_at
	Name         string // Name
	Sub          string
	Com          string // Comment
	Filename     string
	Ext          string
	W            int
	H            int
	Tn_w         int
	Tn_h         int
	Tim          int
	Time         int
	Md5          string
	Fsize        int
	Capcode      string
	Semantic_url string
	Replies      int
	Images       int
	Unique_ips   int
}

func LoadBoards(path string) []JSBoard {
	fData, _ := os.ReadFile(path)

	data := struct {
		Boards []JSBoard
	}{}
	if err := json.Unmarshal(fData, &data); err != nil {
		log.Fatalf("err while unmarshaling (%s)", err)
	}

	return data.Boards
}

type JS struct {
	Page    int
	Threads []JSThread
}

func LoadThreads(path string) chan JSPost {
	out := make(chan JSPost, 100)
	fData, _ := os.ReadFile(path)

	data := []JS{}
	if err := json.Unmarshal(fData, &data); err != nil {
		log.Fatalf("err while unmarshaling (%s)", err)
	}

	var wg sync.WaitGroup

	for _, list := range data {
		for _, t := range list.Threads {
			url := fmt.Sprintf("https://a.4cdn.org/po/thread/%d.json", t.No)
			cPath := fmt.Sprintf("./log/%d.json", t.No)

			wg.Add(1)
			go func(out chan JSPost, url string, cPath string) {
				defer wg.Done()
				if _, err := os.Stat(cPath); !os.IsNotExist(err) {
					//log.Println("got file", cPath)
					for _, p := range LoadPosts(cPath) {
						out <- p
					}
					return
				}

				//log.Println(url)
				res, err := http.Get(url)
				if err != nil {
					log.Printf("Error (%s)\n", err)
					return
				}

				fData, err := io.ReadAll(res.Body)
				res.Body.Close()
				if err != nil {
					log.Printf("err while reading body (%s)\n", err)
					return
				}
				c, _ := os.Create(cPath)
				c.Write(fData)
				defer c.Close()

				pData := struct {
					Posts []JSPost
				}{}
				if err := json.Unmarshal(fData, &pData); err != nil {
					log.Printf("err while unmarshaling (%s)\n", err)
					return
				}
				for _, p := range LoadPosts(cPath) {
					out <- p
				}
			}(out, url, cPath)
		}
	}

	go func(out chan JSPost) {
		wg.Wait()
		out <- JSPost{}
	}(out)

	return out
}

func LoadPosts(path string) []JSPost {
	fData, _ := os.ReadFile(path)

	data := struct {
		Posts []JSPost
	}{}
	if err := json.Unmarshal(fData, &data); err != nil {
		log.Fatalf("err while unmarshaling (%s)", err)
	}

	return data.Posts
}
