package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"sync"

	//"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
)

var count int
var mu sync.Mutex

func cargar(ruta string) {

	abrirCarpeta(ruta)
	//abrirArchivo(ruta)
	//abrirCarpeta("D:/pruebaSoftware/enron_mail_20110402", 0)
	//abrirCarpeta("D:/pruebaSoftware/enron_mail_20110402")
	//abrirArchivo("D:/pruebaSoftware/back/prueba.txt")

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

	for _, fileInfo := range fileInfos {

		if !fileInfo.IsDir() {
			// abrir archivo
			abrirArchivo(ruta + "/" + fileInfo.Name())
		} else {
			count++
			if count > 20 {
				go abrirCarpeta(ruta + "/" + fileInfo.Name())
			} else {
				abrirCarpeta(ruta + "/" + fileInfo.Name())
			}

		}

	}
	log.Println(count)

}

// Estructura de los emails
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

	//Abrir archivo
	f, err := os.Open(ruta)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	email := Email{}

	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	val := reflect.ValueOf(&email).Elem()

	//Leer datos del archivo
	for i := 0; i < val.NumField(); i++ {

		scanner.Scan()

		valueField := val.Field(i)
		parts := strings.Split(scanner.Text(), ":")
		rest := strings.Join(parts[1:], ":")
		valueField.SetString(rest)

	}

	// Leer contenido del email
	bodyContent := ""
	for scanner.Scan() {
		bodyContent += scanner.Text()
	}
	valueField := val.Field(val.NumField() - 1)
	valueField.SetString(bodyContent)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	userJSON, err := json.MarshalIndent(email, "", "  ")

	mu.Lock()
	//Peticion post a ZincSearch
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
	mu.Unlock()
}
