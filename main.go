package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	API_KEY := os.Getenv("RIOT_API_KEY")
	url := "https://euw1.api.riotgames.com/tft/league/v1/grandmaster"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept-Language", "en-GB,en;q=0.5")
	req.Header.Set("Accept-Charset", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("X-Riot-Token", API_KEY)

	res, _ := http.DefaultClient.Do(req)
	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))
}
