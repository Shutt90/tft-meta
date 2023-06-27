package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type ApiKey struct {
	string
}

type Players struct {
	SummonerID   string `json:"summonerId"`
	SummonerName string `json:"summonerName"`
	LeaguePoints uint64 `json:"leaguePoints"`
	Wins         uint64 `json:"wins"`
	Losses       uint64 `json:"losses"`
	Veteran      bool   `json:"veteran"`
	PlayerInfo   playerInfo
}

type Tier struct {
	Tier     string    `json:"tier"`
	LeagueID string    `json:"leagueId"`
	Queue    string    `json:"queue"`
	Name     string    `json:"name"`
	Players  []Players `json:"entries"`
}

type playerInfo struct {
	Id            string `json:"id"`
	AccountId     string `json:"accountId"`
	ProfileIconId uint64 `json:"profileIconId"`
	RevisionDate  uint64 `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
	PUUID         string `json:"puuid"`
	Matches       []string
}

func main() {
	godotenv.Load()
	var apiKey ApiKey
	apiKey.string = os.Getenv("RIOT_API_KEY")
	url := "https://euw1.api.riotgames.com/tft/league/v1/challenger"
	body, err := apiKey.execute("GET", url)
	if err != nil {
		panic(err)
	}

	challenger := Tier{}

	json.Unmarshal(body, &challenger)

	for i, player := range challenger.Players {
		url := "https://euw1.api.riotgames.com/tft/summoner/v1/summoners/by-name/" + player.SummonerName
		body, err = apiKey.execute("GET", url)
		if err != nil {
			panic(err)
		}

		thisPlayerInfo := playerInfo{}

		err = json.Unmarshal(body, &thisPlayerInfo)
		if err != nil {
			fmt.Printf("Could not unmarshal PUUID: %v", err)
		}

		thisPlayerInfo.Matches = apiKey.ScanPlayerMatches(thisPlayerInfo.PUUID)

		challenger.Players[i].PlayerInfo = thisPlayerInfo

	}

	json, err := json.Marshal(challenger)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(json))
}

func (api *ApiKey) execute(method string, url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	headers := map[string][]string{
		"User-Agent":      {"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		"Accept-Language": {"en-GB,en;q=0.5"},
		"Accept-Charset":  {"application/x-www-form-urlencoded; charset=UTF-8"},
		"Origin":          {"https://developer.riotgames.com"},
		"X-Riot-Token":    {api.string},
	}
	req.Header = headers

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (api *ApiKey) ScanPlayerMatches(p string) []string {
	url := fmt.Sprintf("https://europe.api.riotgames.com/tft/match/v1/matches/by-puuid/%s/ids?start=0&count=20", p)
	body, err := api.execute("GET", url)
	if err != nil {
		panic(err)
	}

	var matches []string
	json.Unmarshal(body, &matches)
	if err != nil {
		panic(err)
	}

	return matches
}
