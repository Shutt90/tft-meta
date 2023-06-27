package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Players struct {
	SummonerID   string `json:"summonerId"`
	SummonerName string `json:"summonerName"`
	LeaguePoints uint64 `json:"leaguePoints"`
	Wins         uint64 `json:"wins"`
	Losses       uint64 `json:"losses"`
	Veteran      bool   `json:"veteran"`
}

type Tier struct {
	Tier     string    `json:"tier"`
	LeagueID string    `json:"leagueId"`
	Queue    string    `json:"queue"`
	Name     string    `json:"name"`
	Players  []Players `json:"entries"`
}

func main() {
	godotenv.Load()
	API_KEY := os.Getenv("RIOT_API_KEY")
	url := "https://euw1.api.riotgames.com/tft/league/v1/challenger"
	req, _ := http.NewRequest("GET", url, nil)
	headers := map[string][]string{
		"User-Agent":      {"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		"Accept-Language": {"en-GB,en;q=0.5"},
		"Accept-Charset":  {"application/x-www-form-urlencoded; charset=UTF-8"},
		"Origin":          {"https://developer.riotgames.com"},
		"X-Riot-Token":    {API_KEY},
	}
	req.Header = headers

	res, _ := http.DefaultClient.Do(req)
	body, _ := io.ReadAll(res.Body)

	challenger := Tier{}

	json.Unmarshal(body, &challenger)

	for _, player := range challenger.Players {
		fmt.Println(player.SummonerName)
		fmt.Println(player.SummonerID)
	}
}
