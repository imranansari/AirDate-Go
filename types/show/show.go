package show

import "TvShow/types/episode"

type Show struct {
	Id    int
	Name  string
	Links struct {
		PreviousEpisode struct {
			Href string
		}
		NextEpisode struct {
			Href string
		}
	} `json:"_links, omitempty"`
	Episode   episode.Episode
	Externals struct {
		Tvrage  string
		Thetvdb string
		Imdb    string
	}
}
