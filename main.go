package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"html/template"
	"math/rand"
	"net/http"
)

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

var (
	windV  int
	waterV int
	stat   string
	col    string
)

func inBetween(i, min, max int) bool {
	if (i >= min) && (i <= max) {
		return true
	} else {
		return false
	}
}

func statusUpdate() {
	min := 1
	max := 100
	for {
		rand.Seed(time.Now().UnixNano())
		windV = rand.Intn(max-min+1) + min
		waterV = rand.Intn(max-min+1) + min
		status := &Status{
			Water: waterV,
			Wind:  windV,
		}
		data, _ := json.Marshal(status)
		fmt.Println(string(data))
		if status.Water > 8 || status.Wind > 15 {
			fmt.Printf("Status = BAHAYA \n\n")
			stat = "BAHAYA"
			col = "RED"
		} else if inBetween(status.Water, 6, 8) || inBetween(status.Wind, 7, 15) {
			fmt.Printf("Status = SIAGA \n\n")
			stat = "SIAGA"
			col = "YELLOW"
		} else if status.Water <= 5 && status.Wind <= 6 {
			fmt.Printf("Status = AMAN \n\n")
			stat = "AMAN"
			col = "GREEN"
		}
		time.Sleep(15 * time.Second)
	}
}

func main() {
	go statusUpdate()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var data = map[string]string{
			"Wind":   strconv.Itoa(windV),
			"Water":  strconv.Itoa(waterV),
			"Status": stat,
			"Color":  col,
		}
		var t, err = template.ParseFiles("template.html")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		t.Execute(w, data)
	})
	fmt.Println("starting web server at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)

}
