package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

var file *os.File
var err error
var counter = 1

type requestJson struct {
	Name       string
	SecondName string
	Age        uint8
	IsAdmin    bool
}

func main() {

	date := (time.Now().Format(time.DateTime))

	if _, err := os.Stat("test.txt"); err == nil {
		fmt.Println("Файл уже существует")
		file, _ = os.OpenFile("test.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	} else if errors.Is(err, os.ErrNotExist) {
		file, err = os.Create("test.txt")
		if err != nil {
			fmt.Println("Невозможно создать файл:", err)
			os.Exit(1)
		}
		fmt.Println("Файл успешно создан")
		defer file.Close()
	} else {
		fmt.Println("Возникла непредвиденная ошибка")
	}

	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		file, err = os.OpenFile("test.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Невозможно создать файл:", err)
			os.Exit(1)
		}
		fmt.Println("Файл успешно создан")
		defer file.Close()
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		fmt.Println("Body:", string(body), ",", "Date:", date)
		if err != nil {
			fmt.Println(err)
		}
		file, err = os.OpenFile("test.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		file.Write(([]byte)(`Message: "`))
		file.Write(body)
		file.Write(([]byte)(`" ,`))
		file.Write(([]byte)(date))
		file.Write(([]byte)("\n"))
	})

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat("test.txt"); errors.Is(err, os.ErrNotExist) {
			fmt.Println("Файл не найден")
		} else {
			file.Close()
			os.Remove("test.txt")
			fmt.Println("Файл удален")
		}
	})

	http.HandleFunc("/isExist", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat("test.txt"); err == nil {
			w.Write(([]byte)("true"))
		} else if errors.Is(err, os.ErrNotExist) {
			w.Write(([]byte)("false"))
		} else {
			fmt.Println("Возникла непредвиденная ошибка")
		}
	})

	http.HandleFunc("/requestJson", func(w http.ResponseWriter, r *http.Request) {
		path := "person%s.txt"
		numberPath := fmt.Sprintf(path, strconv.Itoa(counter))
		fileJson, err := os.OpenFile(numberPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		unmarshalBody := requestJson{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
		}
		err = json.Unmarshal(body, &unmarshalBody)
		if err != nil {
			fmt.Println(err)
		}
		marshalAge, err := json.Marshal(unmarshalBody.Age)
		marshalBool, err := json.Marshal(unmarshalBody.IsAdmin)

		fileJson.Write(([]byte)("Name: "))
		fileJson.Write(([]byte)(unmarshalBody.Name))
		fileJson.Write(([]byte)(`, `))
		fileJson.Write(([]byte)("SecondName: "))
		fileJson.Write(([]byte)(unmarshalBody.SecondName))
		fileJson.Write(([]byte)(`, `))
		fileJson.Write(([]byte)("Age: "))
		fileJson.Write(marshalAge)
		fileJson.Write(([]byte)(`, `))
		fileJson.Write(([]byte)("IsAdmin: "))
		fileJson.Write(marshalBool)
		fileJson.Write(([]byte)(`, `))
		fileJson.Write(([]byte)(date))
		fileJson.Write(([]byte)("\n"))
		counter++
	})

	http.ListenAndServe(":80", nil)
}
