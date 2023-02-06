# rle-h2h-history
Allows to load head-to-head history of team-vs-team &amp; player-vs-player matchups in Rocket League Esports.

Hi! My name is 29scroller, also known as Purple Turtle. I am trying myself out in coding, particularly in Golang, and this is my first project.
Rocket League Esports is my passion for many years. I enjoy making predictions of future events, and lately I've been wondering of ways to make them more precise and determined, to make measured forecasts based on public data. Simply put, I want to deduct how teams will play in the future by how these teams / players have played in the past.
Unfortunately, any public resources don't provide that functionality (to be able to see history of matches between certain players / teams of players), so here I am writing my own :)

This project is using zsr.octane.gg API to load data of certain players, teams and matches, which are then being compared between each other to find common matches in the past; then counts the value of the matchup (using completeness of teams and date of match) and sums all matches. So, for any 2 teams of 3 players it shows their past history and summary of their scores in those matchups.
(I am actually so nooby in describing projects more in detail, and have no idea of what it should look like. Maybe the code itself and the comments will tell you more)

It is very early in the making, for now it can calculate the rating of past matchups of two existing teams based on data until February 2nd, 2023. I plan to add much more functionality in the future.

Future plans:
  - Automate adding data of new players and recent matches
  - Personalise input - provide a choice between whether to compare players or teams
  - Detail output - group matches by completeness of rosters, sort by date
  - Give out match predictions based on output
  - Create dataset with only necessary info (for now, info of all player's matches weighs 1.5 GB, of which 98% or more are in-game stats, which aren't used at all)
  - Maybe to design an app for better usage? (seems too complicated for me for now, but we'll see)
