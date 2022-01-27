package main

import (
	"hangman/hangman"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type Page struct {
	PageTitle string
	Title     string
}

type Play struct { // struct output data for the "PLay" pages
	Page
	ToFind       string
	WordHole     string
	Attempt      int
	AttemptImage int
	Tried        []string
	Guess        string
	Advice       string
	Result       string
	Restart      string
}

/*
the function return update struct hangman.Game and update the output in the front page if its win or if its loose
*/
func Result(output Play, data hangman.Game, found bool) (Play, hangman.Game) {
	output.Tried = data.Tried
	output.WordHole = string(data.WordHole)
	output.Attempt = data.Attempt
	if found == true {
		output.WordHole = data.WordRoot
		data.Attempt = 0
		output.Advice = ""
		output.Result = "Win"
		output.Restart = "Clique pour recommencer"
	} else if found == false && data.Attempt == 0 {
		output.WordHole = data.WordRoot
		output.Advice = ""
		output.Result = "Loose"
		output.Restart = "Clic pour recommencer"
	}
	return output, data
}

/*
fonction that return the game data with the input of the user. It return update struct hangman.Game
and play and the booleen (true for word found and false for not found)
*/
func InputWeb(output Play, data hangman.Game) (Play, hangman.Game, bool) {
	var winChar, found, tested bool
	if output.Attempt >= 0 && output.Guess != "" {
		if hangman.IsAlpha(output.Guess) == false {
			output.Advice = "Entrez uniquement des lettres ou des mots s'il vous plait"
		}
		data.Tried, tested = hangman.Tested(data)
		if tested == false {
			output.Advice = "Vous avez déjà testé ceci, essayez encore.."
		}
		if len(data.Input) == 1 && tested == true { // if it's a letter
			winChar, found = hangman.InputChar(data.Input, data.WordRoot, data.WordHole)
			data.Attempt = hangman.Counter(data.Attempt, winChar, 1)
		} else if len(output.Guess) == len(data.WordRoot)-1 && tested == true {
			found = hangman.InputWord(data.Input, output.ToFind)
			data.Attempt = hangman.Counter(data.Attempt, found, 2)
		}
	}
	return output, data, found
}

/*
fonction qui sert la page layout
*/
func viewHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{"Hangman", "Welcome to Hangman Game"}

	tmpl := template.Must(template.ParseFiles("tmpl/layout.gohtml"))

	err := tmpl.ExecuteTemplate(w, "layout", p)
	if err != nil {
		log.Fatalf("Template execution: %s", err)
	}
}

/*
function to serve the page play
*/
func PlayGame1(w http.ResponseWriter, r *http.Request) {
	var found bool
	data, level := hangman.Load()

	if level != 1 || data.Attempt == 0 || data.WordRoot == string(data.WordHole) {
		data, level = hangman.Init(1)
	}
	pageStruct := Page{"Hangman", "Welcome to Hangman Game"}
	outputWeb := Play{pageStruct, data.WordRoot, string(data.WordHole), data.Attempt, data.Attempt, data.Tried, r.FormValue("Guess"), "Propose une lettre ou un mot", "", ""}
	data.Input = r.FormValue("Guess")
	data.Input = strings.ToUpper(r.FormValue("Guess"))

	outputWeb, data, found = InputWeb(outputWeb, data)
	outputWeb, data = Result(outputWeb, data, found)
	if found == true {
		outputWeb.AttemptImage = 11
	} else if found == false && data.Attempt == 0 {
		outputWeb.AttemptImage = 12
	}
	hangman.Save(data, level)

	tmpl := template.Must(template.ParseFiles("tmpl/play.gohtml"))
	err := tmpl.ExecuteTemplate(w, "play", outputWeb)
	if err != nil {
		log.Fatalf("Template execution: %s", err)
	}
	return
}

/*
function to serve the page play2
*/
func PlayGame2(w http.ResponseWriter, r *http.Request) {
	var found bool
	data, level := hangman.Load()
	if level != 2 || data.Attempt == 0 || data.Attempt == 12 || data.Attempt == 11 {
		data, level = hangman.Init(2)
	}

	pageStruct := Page{"Hangman", "Welcome to Hangman Game"}
	outputWeb := Play{pageStruct, data.WordRoot, string(data.WordHole), data.Attempt, data.Attempt, data.Tried, r.FormValue("Guess"), "Propose une lettre ou un mot", "", ""}
	data.Input = r.FormValue("Guess")
	data.Input = strings.ToUpper(r.FormValue("Guess"))

	outputWeb, data, found = InputWeb(outputWeb, data)
	outputWeb, data = Result(outputWeb, data, found)
	if found == true {
		outputWeb.AttemptImage = 11
	} else if found == false && data.Attempt == 0 {
		outputWeb.AttemptImage = 12
	}

	hangman.Save(data, level)

	tmpl := template.Must(template.ParseFiles("tmpl/play2.gohtml"))
	err := tmpl.ExecuteTemplate(w, "play2", outputWeb)
	if err != nil {
		log.Fatalf("Template execution: %s", err)
	}
	return
}

/*
function to serve the page play3
*/
func PlayGame3(w http.ResponseWriter, r *http.Request) {
	var found bool
	data, level := hangman.Load()
	if level != 3 || data.Attempt == 0 {
		data, level = hangman.Init(3)
	}

	pageStruct := Page{"Hangman", "Welcome to Hangman Game"}
	outputWeb := Play{pageStruct, data.WordRoot, string(data.WordHole), data.Attempt, data.Attempt, data.Tried, r.FormValue("Guess"), "Propose une lettre ou un mot", "", ""}
	data.Input = r.FormValue("Guess")
	data.Input = strings.ToUpper(r.FormValue("Guess"))

	outputWeb, data, found = InputWeb(outputWeb, data)
	outputWeb, data = Result(outputWeb, data, found)
	if found == true {
		outputWeb.AttemptImage = 11

	} else if found == false && data.Attempt == 0 {
		outputWeb.AttemptImage = 12
	}
	hangman.Save(data, level)

	tmpl := template.Must(template.ParseFiles("tmpl/play3.gohtml"))
	err := tmpl.ExecuteTemplate(w, "play3", outputWeb)
	if err != nil {
		log.Fatalf("Template execution: %s", err)
	}
	return
}

func main() {
	server := http.NewServeMux()
	server.HandleFunc("/", viewHandler)
	server.HandleFunc("/play", PlayGame1)
	server.HandleFunc("/play2", PlayGame2)
	server.HandleFunc("/play3", PlayGame3)

	server.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	err := http.ListenAndServe(":3000", server)
	if err != nil {
		return
	}
}
