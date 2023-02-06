package headtohead_history

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// CreateClient creates http client to perform API requests and handle responses
func CreateClient() http.Client {
	playerClient := http.Client{}
	fmt.Print("Created client... ")
	return playerClient
}

// CreateRequest creates http request from url string
func CreateRequest(url string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("Created request... ")
	return req
}

// SetHeader sets header for zsr.octane.gg API
func SetHeader(req http.Request) {
	req.Header.Set("Content-Type", "application/json")
	fmt.Print("Set header... ")
}

// GetResponse gets http response from http request using http client
func GetResponse(req *http.Request, clt http.Client) *http.Response {
	res, getErr := clt.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	fmt.Print("Got a response... ")
	return res
}

// CloseBody closes http response's body if it's not empty
func CloseBody(res http.Response) {
	if res.Body != nil {
		defer res.Body.Close()
		fmt.Print("Closed body... ")
	} else {
		fmt.Println("Body is nil, something went wrong.")
	}
}

// WriteToByteSlice writes an http response into byte slice
func WriteToByteSlice(res http.Response) []byte {
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	fmt.Print("Wrote info in byte slice... ")
	return body
}

// UrlToByteSlice combines all the previous steps for converting URL to byte slice with data.
func UrlToByteSlice(url string) (rawData []byte) {
	client := CreateClient()
	request := CreateRequest(url)
	SetHeader(*request)
	response := GetResponse(request, client)
	rawData = WriteToByteSlice(*response)
	CloseBody(*response)
	fmt.Println()
	return
}

// WriteMatchesOfPlayerToFile writes all found matches of player to file named as player's slug
func WriteMatchesOfPlayerToFile(player Player) {
	f, err := os.Create("C:\\Users\\29scroller\\Downloads\\TEAMVSTEAM\\matches_of_players\\" + player.Slug + ".txt")
	if err != nil {
		fmt.Println("Could not create file :(")
		panic(err)
	}
	page := 1
	for {
		url := "https://zsr.octane.gg/matches?mode=3&page=" + fmt.Sprint(page) + "&player=" + player.Id
		rawData := UrlToByteSlice(url)
		rawData = []byte(strings.TrimSpace(string(rawData)))
		endString := `{"matches":[],"page":` + fmt.Sprint(page) + `,"perPage":100,"pageSize":0}`
		if string(rawData) == endString {
			f.Close()
			break
		}
		f.Write(rawData)
		f.WriteString("\n")
		fmt.Printf("Wrote page %d of player %s \n", page, player.Slug)
		page += 1
	}
	fmt.Println("Wrote info to file", player.Slug)
}

// LoadPlayerMatchesFromFile loads info wrote by WriteMatchesOfPlayerToFile from file associated with player's slug
func LoadPlayerMatchesFromFile(player Player) []Match {
	rawData, err := os.ReadFile("C:\\Users\\29scroller\\Downloads\\TEAMVSTEAM\\matches_of_players\\" + player.Slug + ".txt")
	if err != nil {
		panic(err)
	}
	stringRawData := strings.TrimSpace(string(rawData))
	splitRawData := regexp.MustCompile("\r?\n").Split(stringRawData, -1)
	playerMatches := make([]MultipleMatches, len(splitRawData))
	var matchesOnly []Match
	for i := 0; i < len(splitRawData); i++ {
		tempByteSlice := []byte(splitRawData[i])
		UnmarshalObject(tempByteSlice, &playerMatches[i])
		matchesOnly = append(matchesOnly, playerMatches[i].Matches...)
	}
	fmt.Printf("Matches count of %s = %d\n", player.Tag, len(matchesOnly))
	return matchesOnly
}

// UnmarshalObject unmarshals byte slice into any struct.
func UnmarshalObject[T any](body []byte, object T) {
	jsonErr := json.Unmarshal(body, &object)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}

//UserEnteringTeam handles dialogue with user for entering team name.
func UserEnteringTeam(numberAdj, teamName string) (teamInfo Team, teamFound bool) {
	fmt.Println("Write", numberAdj, "team's name")
	reader := bufio.NewReader(os.Stdin)
	var err error
	teamName, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Finding ID of", teamName)
	teamInfo, teamFound = FindTeamByName(teamName)
	fmt.Printf("Found team = %t, team ID is %s \n", teamFound, teamInfo.Id)
	return
}
