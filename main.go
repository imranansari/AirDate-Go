package main

import (
	"TvShow/icon"
	. "TvShow/types/episode"
	. "TvShow/types/show"
	"encoding/json"
	"fmt"
	"github.com/getlantern/systray"
	"gopkg.in/toast.v1"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

var imdbLookup = "http://api.tvmaze.com/lookup/shows?imdb="

func main() {
	systray.Run(onReady, onExit)

	// reader := bufio.NewReader(os.Stdin)
	// reader.ReadLine()
}

func setShowValues(url string, show *Show) {
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &show)

	if len(show.Links.NextEpisode.Href) > 0 {
		setShowEpisode(show.Links.NextEpisode.Href, show)
	}
}

func setShowEpisode(url string, show *Show) {
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var episode Episode

	json.Unmarshal(body, &episode)
	show.Episode = episode
}

func onReady() {
	load()

	systray.SetIcon(icon.Data)
	systray.SetTitle("Next Aired")
	systray.SetTooltip("Next Aired")

	mReload := systray.AddMenuItem("Reload", "Check again!")
	systray.AddSeparator()
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		for {
			select {
			case <-mQuitOrig.ClickedCh:
				systray.Quit()
			case <-mReload.ClickedCh:
				load()
			}
		}
	}()
}

func onExit() {
	// clean up here
}

func load() {
	var ids []string
	ids = append(ids, "tt7371666")
	ids = append(ids, "tt1844624")
	ids = append(ids, "tt0460681")
	ids = append(ids, "tt6763664")
	ids = append(ids, "tt7008682")
	ids = append(ids, "tt0944947")
	ids = append(ids, "tt7493974")
	ids = append(ids, "tt3006802")
	ids = append(ids, "tt3107288")

	var shows []Show

	for _, id := range ids {
		var show Show
		setShowValues(imdbLookup+id, &show)
		shows = append(shows, show)
	}

	for _, v := range shows {

		var title, message string
		send := false

		title = fmt.Sprintf("%s\n", v.Name)

		if len(v.Links.NextEpisode.Href) > 0 {
			if v.Episode.DaysLeft() <= 2 {
				send = true
			}
			message = fmt.Sprintf("%s (S%dE%d)\nTime Left: %s\n\n",
				v.Episode.Name,
				v.Episode.Season,
				v.Episode.Number,
				v.Episode.TimeLeft(),
			)
		} else {
			message = fmt.Sprintln("Next Episode not known\n")
		}

		if send {
			err := alert(title, message, "assets/main.png", "https://www.imdb.com/title/"+v.Externals.Imdb+"/episodes", v.Episode.Url)
			if err != nil {
				panic(err)
			}
		}
	}
}

func toastNotification(title, message, appIcon string) toast.Notification {
	// NOTE: a real appID is required since Windows 10 Fall Creator's Update,
	// issue https://github.com/go-toast/toast/issues/9
	appID := "{1AC14E77-02E7-4E5D-B744-2EB1AE5198B7}\\WindowsPowerShell\\v1.0\\powershell.exe"
	return toast.Notification{
		AppID:   appID,
		Title:   title,
		Message: message,
		Icon:    appIcon,
	}
}

// Alert displays a desktop notification and plays a default system sound.
func alert(title, message, appIcon, imdb, tvmaze string) error {
	var err error
	iconPath := ""
	if appIcon != "" {
		iconPath, err = filepath.Abs(appIcon)
		if err != nil {
			return err
		}
	}
	note := toastNotification(title, message, iconPath)
	note.Actions = []toast.Action{
		{"protocol", "IMDB", imdb},
		{"protocol", "TVMAZE", tvmaze},
	}
	note.Audio = toast.Default
	return note.Push()
}
