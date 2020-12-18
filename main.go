package main

import (
	"encoding/json"
	"fmt"
	"github.com/code7unner/student-parser/models"
	"net/http"
	"os"
	"sync"
)

const (
	bufSize = 1000000
	workers = 10
)

var accessToken string

func init() {
	accessToken = os.Getenv("ACCESS_TOKEN")
	if accessToken == "" {
		panic("ACCESS_TOKEN is required")
	}
}

func main() {
	groupIDs := []string{"math_100", "inf_bu", "ege_matn", "physics_100", "ege", "egeoge_math"}

	vkID := make(chan int, bufSize)
	for _, id := range groupIDs {
		offset := 0
		go getVkIDs(vkID, id, offset)
	}

	wg := sync.WaitGroup{}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			for {
				select {
				case id, more := <-vkID:
					if !more {
						wg.Done()
						return
					}

					getVkProfile(id)
					getVkGroups(id)
				}
			}
		}()
	}
	wg.Wait()
}

func getVkGroups(id int) {
	url := "https://api.vk.com/method/users.getSubscriptions?user_id=%d&extended=1&count=200&v=5.126&access_token=%s"

	url = fmt.Sprintf(url, id, accessToken)
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	var data interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}
	resp.Body.Close()

	fmt.Println(data)
}

func getVkProfile(id int) {
	url := "https://api.vk.com/method/users.get?user_id=%d&fields=verified,interests,sex&extended=1&count=1000&v=5.126&access_token=%s"

	url = fmt.Sprintf(url, id, accessToken)
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	var data models.VkUserModel
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}
	resp.Body.Close()

	fmt.Println(data)
}

func getVkIDs(vkID chan int, id string, offset int) {
	for {
		url := "https://api.vk.com/method/groups.getMembers?group_id=%s&sort=id_asc&offset=%d&count=1000&v=5.126&access_token=%s"

		url = fmt.Sprintf(url, id, offset, accessToken)
		resp, err := http.Get(url)
		if err != nil {
			break
		}

		offset += 1000

		var data models.VkIDModel
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			break
		}
		resp.Body.Close()

		if len(data.Response.Items) == 0 {
			break
		}

		for _, item := range data.Response.Items {
			vkID <- item
		}
	}

	close(vkID)
}
