package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"unicode"
)

func cargar() {

	abrirCarpeta("D:/pruebaSoftware/enron_mail_20110402")
	//abrirArchivo("D:/pruebaSoftware/back/prueba.txt")
	//http.HandleFunc("/debug/pprof/", pprof.Index)
	//http.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	//http.HandleFunc("/debug/pprof/profile", pprof.Profile)
	//http.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
}

func abrirCarpeta(ruta string) {

	// Abrir la carpeta
	dir, err := os.Open(ruta)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dir.Close()

	// Leer la lista de archivos en la carpeta
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Imprimir el nombre de cada archivo
	for _, fileInfo := range fileInfos {
		//fmt.Println(fileInfo.Name())
		iter := true
		for _, c := range fileInfo.Name() {

			if !unicode.IsDigit(c) {
				iter = false
			}
		}
		//time.Sleep(100)

		if iter == true {
			// abrir archivo
			abrirArchivo(ruta + "/" + fileInfo.Name())
		} else {
			 abrirCarpeta(ruta + "/" + fileInfo.Name())
		}

	}

}

type Email struct {
	Message_ID                string
	Date                      string
	From                      string
	To                        string
	Subject                   string
	Mime_Version              string
	Content_Type              string
	Content_Transfer_Encoding string
	X_From                    string
	X_To                      string
	X_cc                      string
	X_bcc                     string
	X_Folder                  string
	X_Origin                  string
	X_FileName                string
	Body                      string
}

func abrirArchivo(ruta string) {

	f, err := os.Open(ruta)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	email := Email{}

	scanner := bufio.NewScanner(f)
	val := reflect.ValueOf(&email).Elem()

	for i := 0; i < val.NumField(); i++ {

		scanner.Scan()
		valueField := val.Field(i)
		parts := strings.Split(scanner.Text(), ":")
		rest := strings.Join(parts[1:], ":")
		valueField.SetString(rest)
		//typeField := val.Type().Field(i)
		//fmt.Printf("%s: %v\n", typeField.Name, valueField.Interface())

	}

	bodyContent := ""
	for scanner.Scan() {
		bodyContent += scanner.Text()
	}
	valueField := val.Field(val.NumField() - 1)
	valueField.SetString(bodyContent)

	/* print the email struct
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		fmt.Printf("%s: %v\n", typeField.Name, valueField.Interface())

	}
	*/

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	//println("=======")
	userJSON, err := json.MarshalIndent(email, "", "  ")
	//fmt.Println(string(userJSON))
	//=============

	req, err := http.NewRequest("POST", "http://localhost:4080/api/games3/_doc", strings.NewReader(string(userJSON)))
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
	fmt.Println(string(body))

}
