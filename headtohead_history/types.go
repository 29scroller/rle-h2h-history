package headtohead_history

//This is the data structure for parsing json data from zsr.octane.gg.
//Only data necessary for program to work is defined, stats and everything else is omitted, though stored in database
//(for there's seemingly no way to filter fields from the start, needs research).
//Also added AdjustedResult type to store adjusted results in Games variable.

type MultipleMatches struct {
	Matches []Match `json:"matches"`
}

type Match struct {
	Id       string `json:"_id"`
	OctaneId string `json:"octane_id"`
	MEvent   Event  `json:"event"`
	Date     string `json:"date"`
	Blue     Games  `json:"blue"`
	Orange   Games  `json:"orange"`
}

type Event struct {
	Id   string `json:"_id"`
	Name string `json:"name"`
	Tier string `json:"tier"`
}

type Games struct {
	Score          uint8              `json:"score"`
	Winner         bool               `json:"winner"`
	TeamUp         TeamForMatches     `json:"team"`
	PlayerUp       []PlayerForMatches `json:"players"`
	AdjustedResult AdjustedResult
}

type TeamForMatches struct {
	Team Team `json:"team"`
}

type MultiplePlayers struct {
	Players    []Player `json:"players"`
}

type PlayerForMatches struct {
	Player Player `json:"player"`
}

type Player struct {
	Id         string `json:"_id"`
	Slug       string `json:"slug"`
	Tag        string `json:"tag"`
	PTeam      Team   `json:"team"`
	Teams      []Team `json:"teams"`
	Substitute bool   `json:"substitute"`
	Coach      bool   `json:"coach"`
}

type Team struct {
	Id   string `json:"_id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type AdjustedResult struct {
	AdjustedSeries float64
	AdjustedGames  float64
}
