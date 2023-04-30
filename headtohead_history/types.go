package headtohead_history

import (
	"time"
)

//This is the data structure for parsing json data from zsr.octane.gg.
//Only data necessary for program to work is defined, stats and everything else is omitted, though stored in database
//(for there's seemingly no way to filter fields from the start, needs research).

type MultipleMatches struct {
	Matches []Match `json:"matches"`
}

type Match struct {
	Id                             string `json:"_id"`
	OctaneId                       string `json:"octane_id"`
	MEvent                         Event  `json:"event"`
	Stage                          Stage  `json:"stage"`
	Date                           string `json:"date"`
	Format                         Format `json:"format"`
	Blue                           Games  `json:"blue"`
	Orange                         Games  `json:"orange"`
	DateTime                       time.Time
	BlueT, BlueO, OrangeT, OrangeO uint8
	btoo, otbo                     bool
}

type Event struct {
	Id     string `json:"_id"`
	Name   string `json:"name"`
	Region string `json:"region"`
	Tier   string `json:"tier"`
}

type Stage struct {
	Lan bool `json:"lan"`
}

type Format struct {
	Length int8 `json:"length"`
}

type Games struct {
	Score    uint8              `json:"score"`
	Winner   bool               `json:"winner"`
	TeamUp   TeamForMatches     `json:"team"`
	PlayerUp []PlayerForMatches `json:"players"`
}

type TeamForMatches struct {
	Team Team `json:"team"`
}

type MultiplePlayers struct {
	Players []Player `json:"players"`
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
	Country    string `json:"country"`
	InTeam     bool
}

type Team struct {
	Id    string `json:"_id"`
	Slug  string `json:"slug"`
	Name  string `json:"name"`
	Image string `json:"image"`
}
