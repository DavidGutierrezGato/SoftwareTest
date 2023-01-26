package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

type MyStruct struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		MaxScore float64 `json:"max_score"`
		Hits     []struct {
			Index     string  `json:"_index"`
			Type      string  `json:"_type"`
			ID        string  `json:"_id"`
			Score     float64 `json:"_score"`
			Timestamp string  `json:"@timestamp"`
			Source    struct {
				From    string `json:"From"`
				Subject string `json:"Subject"`
				To      string `json:"To"`
				Body    string `json:"Body"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	r := chi.NewRouter()

	
	cargar()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Get("/find/{find}", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		query := `{
			"search_type": "match",
			"query":
			{
				"term": "something",
				"start_time": "2023-01-24T15:08:48.777Z",
        		"end_time": "2023-11-25T16:08:48.777Z"
				
			},
			"from": 0,
			"max_results": 20,
			"_source": ["From","Subject","To","Body"]
		}`
		newTerm := chi.URLParam(r, "find")
		newTerm2 := `"` + newTerm + `"`
		println("term: " + newTerm2)
		query = strings.Replace(query, `"term": "something"`, `"term":`+newTerm2, -1)
		req, err := http.NewRequest("POST", "http://localhost:4080/api/games3/_search", strings.NewReader(query))
		if err != nil {
			log.Fatal(err)

		}
		req.SetBasicAuth("admin", "Complexpass#123")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		log.Println(resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(string(body))

		var myStruct MyStruct
		json.Unmarshal([]byte(string(body)), &myStruct)
		//fmt.Printf("%+v", myStruct)
		fmt.Printf("---------")
		//fmt.Printf("%+v", myStruct.Hits.Hits)
		//str := fmt.Sprintf("%+v", myStruct.Hits.Hits)

		out, err := json.Marshal(myStruct.Hits.Hits)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(out))

		arreglado := strings.Replace(string(out), "_source", "source", -1)
		w.Write([]byte(arreglado))
	})

	http.ListenAndServe(":3000", r)
}
