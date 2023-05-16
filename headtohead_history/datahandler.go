package headtohead_history

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
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
func WriteMatchesOfPlayerToFile(player Player, rawData []byte) {
	f, err := os.Create("matches_of_players\\" + player.Slug + ".txt")
	if err != nil {
		fmt.Println("Could not create file :(")
		panic(err)
	}
	if len(rawData) == 0 {
		rawData = CollectMatchesOfPlayerInByteSlice(player)
	}
	WriteByteSliceToPlayerFile(*f, player, rawData)
	fmt.Printf("Wrote info to file %s.txt\n", player.Slug)
}

// WriteByteSliceToPlayerFile writes matches data into set file
func WriteByteSliceToPlayerFile(file os.File, player Player, rawData []byte) {
	file.Write(rawData)
	file.Close()
	fmt.Printf("Wrote info of %s to file %s.txt\n", player.Tag, player.Slug)
}

// CheckPlayerForNewMatches indicates whether the player has new matches
func CheckPlayerForNewMatches(player Player) {
	data, sliceChanged := BiggestOfFileAndJsonForPlayer(player)
	if sliceChanged {
		WriteMatchesOfPlayerToFile(player, data)
		fmt.Printf("Wrote updated info in player's file, %s.txt\n", player.Slug)
	} else {
		fmt.Printf("No changes in data for %s\n", player.Tag)
	}
}

// CollectMatchesOfPlayerInByteSlice collects all matches of player in a byte slice
func CollectMatchesOfPlayerInByteSlice(player Player) (totalData []byte) {
	page := 1
	for {
		url := "https://zsr.octane.gg/matches?mode=3&page=" + fmt.Sprint(page) + "&player=" + player.Id
		rawData := UrlToByteSlice(url)
		rawData = []byte(strings.TrimSpace(string(rawData)))
		endString := `{"matches":[],"page":` + fmt.Sprint(page) + `,"perPage":100,"pageSize":0}`
		if string(rawData) == endString {
			return totalData
		}
		totalData = append(totalData, rawData...)
		totalData = append(totalData, "\n"...)
		fmt.Printf("Read page %d of player %s \n", page, player.Slug)
		page += 1
	}
}

// CollectAllMatchesInByteSlice collects all matches with certain parameters in a byte slice
func CollectAllMatchesInByteSlice() (totalData []byte) {
	page := 1
	for {
		url := "https://zsr.octane.gg/matches?mode=3&sort=date%3Aasc&page=" + fmt.Sprint(page)
		rawData := UrlToByteSlice(url)
		rawData = []byte(strings.TrimSpace(string(rawData)))
		endString := `{"matches":[],"page":` + fmt.Sprint(page) + `,"perPage":100,"pageSize":0}`
		if string(rawData) == endString {
			return totalData
		}
		totalData = append(totalData, rawData...)
		totalData = append(totalData, "\n"...)
		fmt.Printf("Read page %d of all matches \n", page)
		page += 1
	}
}

// CollectAllMatchesInCSV collects all matches with certain parameters in matches.csv file
func CollectAllMatchesInCSV() {
	page := 1
	file, err := os.Create("matches.csv")

	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(file)

	w.Write(CreateCSVHeader())
	//file.Close()
	for {
		url := "https://zsr.octane.gg/matches?mode=3&sort=date:asc&page=" + fmt.Sprint(page)
		rawData := UrlToByteSlice(url)
		var tempMatches MultipleMatches
		UnmarshalObject(rawData, &tempMatches)
		data := MatchesToStringSlice(tempMatches)
		if len(data) == 0 {
			break
		}
		file, err = os.Open("matches.csv")
		if err != nil {
			log.Fatalln("failed to open file", err)
		}
		if err := w.WriteAll(data); err != nil {
			log.Fatalln("error writing record to file", err)
		}
		file.Close()
		/* for _, row := range data {
			file, err = os.Open("matches.csv")
			if err != nil {
				log.Fatalln("failed to open file", err)
			}
			if err := w.Write(row); err != nil {
				log.Fatalln("error writing record to file", err)
			}
			file.Close()
		} */
		fmt.Println("Wrote info from page", page)
		page++
		w.Flush()
	}
}

// BiggestOfFileAndJsonForPlayer returns biggest byte slice between old data and new data
func BiggestOfFileAndJsonForPlayer(player Player) (biggestSlice []byte, hasChanged bool) {
	oldData := FileToByteSlice(player)
	newData := CollectMatchesOfPlayerInByteSlice(player)
	if len(newData) > len(oldData) {
		return newData, true
	} else {
		return oldData, false
	}
}

// AppendMatchesOfPlayerToFile writes data to player file
func AppendMatchesOfPlayerToFile(player Player, newData []byte) {
	f, err := os.OpenFile("matches_of_players\\"+player.Slug+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Could not create file :(")
		panic(err)
	}
	defer f.Close()
	_, err2 := f.WriteString(string(newData))
	if err2 != nil {
		fmt.Printf("Could not write new data to file %s.txt\n", player.Slug)
	} else {
		fmt.Printf("New data has been written to %s.txt\n", player.Slug)
	}
}

// LoadPlayerMatchesFromFile loads info wrote by WriteMatchesOfPlayerToFile from file associated with player's slug
func LoadPlayerMatchesFromFile(player Player) (matchesOnly []Match) {
	stringRawData := strings.TrimSpace(string(FileToByteSlice(player)))
	splitRawData := regexp.MustCompile("\r?\n").Split(stringRawData, -1)
	playerMatches := make([]MultipleMatches, len(splitRawData))
	for i := 0; i < len(splitRawData); i++ {
		tempByteSlice := []byte(splitRawData[i])
		UnmarshalObject(tempByteSlice, &playerMatches[i])
		matchesOnly = append(matchesOnly, playerMatches[i].Matches...)
	}
	fmt.Printf("Matches count of %s = %d\n", player.Tag, len(matchesOnly))
	return
}

// DoesPlayerFileExist returns whether player file exists or not
func DoesPlayerFileExist(player Player) bool {
	_, error := os.Stat("matches_of_players\\" + player.Slug + ".txt")
	if os.IsNotExist(error) {
		return false
	} else {
		return true
	}
}

// CheckIfPlayerHasFile checks if player has the file and if not, creates it
func CheckIfPlayerHasFile(player Player) {
	var rawData []byte
	if !DoesPlayerFileExist(player) {
		WriteMatchesOfPlayerToFile(player, rawData)
	}
}

// CheckIfTeamPlayersHaveFiles checks if all players on the team have files
func CheckIfTeamPlayersHaveFiles(players []Player) {
	for i := 0; i < len(players); i++ {
		CheckIfPlayerHasFile(players[i])
	}
}

// FileToByteSlice reads player file and returns byte slice with file's data
func FileToByteSlice(player Player) (rawData []byte) {
	rawData, err := os.ReadFile("matches_of_players\\" + player.Slug + ".txt")
	if err != nil {
		panic(err)
	}
	return
}

// UnmarshalObject unmarshals byte slice into any struct.
func UnmarshalObject[T any](body []byte, object T) {
	jsonErr := json.Unmarshal(body, &object)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}

// UserEnteringTeams handles entering teams by user
func UserEnteringTeams() (firstTeamInfo, secondTeamInfo Team) {
	teamFound := false
	for !teamFound {
		firstTeamInfo, teamFound = UserEnteringTeam("first")
	}
	teamFound = false
	for !teamFound {
		secondTeamInfo, teamFound = UserEnteringTeam("second")
	}
	return
}

// UserEnteringTeam handles dialogue with user for entering team name.
func UserEnteringTeam(numberAdj string) (teamInfo Team, teamFound bool) {
	fmt.Println("Write", numberAdj, "team's name")
	reader := bufio.NewReader(os.Stdin)
	var err error
	teamName, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Finding info of", teamName)
	teamInfo, teamFound = FindTeamByName(teamName)
	fmt.Printf("Found team = %t, \n", teamFound)
	return
}

// UserEnteringAllPlayers handles entering players by user
func UserEnteringAllPlayers() (firstTeamPlayers, secondTeamPlayers []Player) {
	firstTeamPlayers = UserEnteringPlayersOfTeam("first")
	secondTeamPlayers = UserEnteringPlayersOfTeam("second")
	return
}

// UserEnteringPlayersOfTeam handles entering players of one team by user
func UserEnteringPlayersOfTeam(numberAdj string) (teamPlayers []Player) {
	teamPlayers = make([]Player, 3)
	fmt.Println("Write", numberAdj, "team's players")
	playerfound := false
	for !playerfound {
		teamPlayers[0], playerfound = UserEnteringPlayer("first")
	}
	playerfound = false
	for !playerfound {
		teamPlayers[1], playerfound = UserEnteringPlayer("second")
	}
	playerfound = false
	for !playerfound {
		teamPlayers[2], playerfound = UserEnteringPlayer("third")
	}
	return
}

// UserEnteringPlayer handles entering info of one player by user
func UserEnteringPlayer(numberAdj string) (playerInfo Player, playerfound bool) {
	fmt.Println("Write", numberAdj, "player's tag")
	reader := bufio.NewReader(os.Stdin)
	var err error
	playerTag, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Finding info of", playerTag)
	playerInfo, playerfound = FindPlayerByTag(playerTag)
	fmt.Printf("Found team = %t, \n", playerfound)
	return
}

// UserChoosingInput handles choosing input method by user
func UserChoosingInput() (userInput int) {
	fmt.Println("Do you want to input teams or distinct players? Type 1 for teams or 2 for players")
	fmt.Scanln(&userInput)
	return
}

// UserChoosingWhetherToCheckPlayers handles choosing option of checking players' files by user
func UserChoosingWhetherToCheckPlayers() (userInput int) {
	fmt.Println("Do you want to check players' files for updates? Type 1 to check or 0 to not check")
	fmt.Scanln(&userInput)
	return
}

// CompileAllPlayerRecords writes all matches' data into matches.csv file
func CompileAllPlayerRecords() {
	rawData := CollectAllMatchesInByteSlice()
	var matches MultipleMatches
	UnmarshalObject(rawData, &matches)
	f, err := os.Create("matches.csv")

	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(f)
	defer w.Flush()

	w.Write(CreateCSVHeader())
	data := MatchesToStringSlice(matches)

	for _, row := range data {
		f, err = os.Open("matches.csv")
		if err != nil {
			log.Fatalln("failed to open file", err)
		}
		if err := w.Write(row); err != nil {
			log.Fatalln("error writing record to file", err)
		}
		f.Close()
	}
}

// CreateCSVHeader creates the first string of matches.csv file
func CreateCSVHeader() (header []string) {
	header = append(header, "Match ID", "Event ID", "Event Name", "Event Region", "Event Tier")
	header = append(header, "Is match on LAN?", "Match Date", "Best of _")
	header = append(header, "Blue Team ID", "Blue Team Name", "Blue Team Image")
	header = append(header, "Blue Player ID", "Blue Player Tag", "Blue Player Country")
	header = append(header, "Blue Games Score", "Blue Wins Series?")
	header = append(header, "Orange Team ID", "Orange Team Name", "Orange Team Image")
	header = append(header, "Orange Player ID", "Orange Player Tag", "Orange Player Country")
	header = append(header, "Orange Games Score", "Orange Wins Series?")
	return
}

// MatchesToStringSlice splits MultipleMatches struct into slice of slices of strings
func MatchesToStringSlice(matches MultipleMatches) (matchSlice [][]string) {
	for i := 0; i < len(matches.Matches); i++ {
		matchSlice = append(matchSlice, FillStringSliceWithAllData(matches.Matches[i])...)
	}
	return
}

// FillStringSliceWithMatchData fills slice of strings with match data
func FillStringSliceWithMatchData(match Match) (stringSlice []string) {
	stringSlice = append(stringSlice, match.Id, match.MEvent.Id, match.MEvent.Name, match.MEvent.Region, match.MEvent.Tier)
	stringSlice = append(stringSlice, strconv.FormatBool(match.Stage.Lan), match.Date, strconv.FormatInt(int64(match.Format.Length), 10))
	stringSlice = append(stringSlice, match.Blue.TeamUp.Team.Id, match.Blue.TeamUp.Team.Name, match.Blue.TeamUp.Team.Image)
	stringSlice = append(stringSlice, "", "", "")
	stringSlice = append(stringSlice, strconv.FormatUint(uint64(match.Blue.Score), 10), strconv.FormatBool(match.Blue.Winner))
	stringSlice = append(stringSlice, match.Orange.TeamUp.Team.Id, match.Orange.TeamUp.Team.Name, match.Orange.TeamUp.Team.Image)
	stringSlice = append(stringSlice, "", "", "")
	stringSlice = append(stringSlice, strconv.FormatUint(uint64(match.Orange.Score), 10), strconv.FormatBool(match.Orange.Winner))
	return
}

// FillStringSliceWithAllData fills slice of slices of strings with match data and players' data
func FillStringSliceWithAllData(match Match) (matchSlice [][]string) {
	var tempSlice []string
	for i := 0; i < len(match.Blue.PlayerUp); i++ {
		for j := 0; j < len(match.Orange.PlayerUp); j++ {
			tempSlice = FillStringSliceWithMatchData(match)
			tempSlice[11] = match.Blue.PlayerUp[i].Player.Id
			tempSlice[12] = match.Blue.PlayerUp[i].Player.Tag
			tempSlice[13] = match.Blue.PlayerUp[i].Player.Country
			tempSlice[19] = match.Orange.PlayerUp[j].Player.Id
			tempSlice[20] = match.Orange.PlayerUp[j].Player.Tag
			tempSlice[21] = match.Orange.PlayerUp[j].Player.Country
			matchSlice = append(matchSlice, tempSlice)
		}
	}
	return matchSlice
}
