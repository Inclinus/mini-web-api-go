package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"unicode"
)

func printLnInNav(msg string, w *http.ResponseWriter) {
	_, err := fmt.Fprintln(*w, msg)
	if err != nil {
		fmt.Println("Une erreur est survenue.")
		return
	}
}

func printInNav(msg string, w *http.ResponseWriter) {
	_, err := fmt.Fprint(*w, msg)
	if err != nil {
		fmt.Println("Une erreur est survenue.")
		return
	}
}

func getTime() string {
	timeNow := time.Now()
	hour := timeNow.Hour()
	min := timeNow.Minute()

	return fmt.Sprintf("%02dh%02d", hour, min)
}

func getDice() int {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	return random.Intn(1000) + 1
}

func getDices(forcedDice *string, w *http.ResponseWriter) {
	var dices = []int{2, 4, 6, 8, 10, 12, 20, 100}
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	diceCorrespondance := map[string]int{"d2": 2, "d4": 4, "d6": 6, "d8": 8, "d10": 10, "d12": 12, "d20": 20, "d100": 100}
	var integers []int
	if forcedDice != nil {
		for i := 0; i < 15; i++ {
			dice := diceCorrespondance[*forcedDice]
			if dice != 2 && dice != 4 && dice != 6 && dice != 8 && dice != 10 && dice != 12 && dice != 20 && dice != 100 {
				printLnInNav("ERROR : Merci de préciser un type correct", w)
				return
			}
			integers = append(integers, random.Intn(dice)+1)
		}
	} else {
		for i := 0; i < 15; i++ {
			randomDice := dices[random.Intn(8)]
			integers = append(integers, random.Intn(randomDice)+1)
		}
	}
	for _, i := range integers {
		if i > 20 {
			printInNav(fmt.Sprintf("%03d ", i), w)
		} else {
			printInNav(fmt.Sprintf("%d ", i), w)
		}

	}
}

func randomizeWords(sentence string) string {
	str := strings.Fields(sentence)

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	random.Shuffle(len(str), func(a, b int) {
		str[a], str[b] = str[b], str[a]
	})

	return strings.Join(str, " ")
}

func semiCapitalizeSentence(sentence string) string {
	runeSentence := []rune(sentence)

	for i := 0; i < len(runeSentence); i++ {
		if i%2 == 0 {
			runeSentence[i] = unicode.ToUpper(runeSentence[i])
		}
	}

	return string(runeSentence)
}

func timeHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		printLnInNav(getTime(), &w)
	}
}

func diceHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		printInNav(fmt.Sprintf("%04d", getDice()), &w)
	}
}

func dicesHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		dice := req.URL.Query().Get("type")
		if dice == "" {
			getDices(nil, &w)
		} else {
			getDices(&dice, &w)
		}
	}
}

func randomizeWordHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		if err := req.ParseForm(); err != nil {
			fmt.Println("Something went bad")
			printLnInNav("Something went bad", &w)
			return
		}
		for key, value := range req.PostForm {
			if key != "words" {
				printLnInNav("Veuillez rentrer uniquement le paramètre words.", &w)
				return
			} else {
				if value[0] != "" {
					printLnInNav(randomizeWords(value[0]), &w)
				} else {
					printLnInNav("Veuillez rentrer une phrase en paramètre.", &w)
					return
				}
			}
		}
	}
}

func semiCapitalizeSentenceHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		if err := req.ParseForm(); err != nil {
			fmt.Println("Something went bad")
			printLnInNav("Something went bad", &w)
			return
		}
		for key, value := range req.PostForm {
			if key != "sentence" {
				printLnInNav("Veuillez rentrer uniquement le paramètre sentence.", &w)
				return
			} else {
				if value[0] != "" {
					printLnInNav(semiCapitalizeSentence(value[0]), &w)
				} else {
					printLnInNav("Veuillez rentrer une phrase en paramètre.", &w)
					return
				}
			}
		}
	}
}

func main() {
	http.HandleFunc("/", timeHandler)
	http.HandleFunc("/dice", diceHandler)
	http.HandleFunc("/dices", dicesHandler)
	http.HandleFunc("/randomize-words", randomizeWordHandler)
	http.HandleFunc("/semi-capitalize-sentence", semiCapitalizeSentenceHandler)

	err := http.ListenAndServe(":4567", nil)
	if err != nil {
		fmt.Printf("ERROR OCCURRED WHILE LAUNCHING SERVER :\n%v", err)
		return
	}
}
