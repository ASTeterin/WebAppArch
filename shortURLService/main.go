package main

import (
    "encoding/json"
    "fmt"      // пакет для форматированного ввода вывода
    "log"      // пакет для логирования
    "net/http" // пакет для поддержки HTTP протокола
)

type PathStr map[string]string

type PathConf struct {
    Paths struct {
        p PathStr
    } `json:"paths"`
}



func ParsePath() {

    //return paths
}

func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {
    //var paths := ParsePath()
    var pathConfiguration interface{}
    //var con PathConf
    jsonData := []byte(`
    {
        "paths" : {
            "/go-http": "https://golang.org/pkg/net/http/",
            "/go-gophers" : "https://github.com/shalakhin/gophericons/blob/master/preview.jpg"
        }
    }`)

    err := json.Unmarshal(jsonData, &pathConfiguration)
    if err != nil {
        log.Println(err)
    }

    //var str = con.Paths.p["/go-http"]
    var str = fmt.Sprintf("%v", pathConfiguration)
    fmt.Println(pathConfiguration)
    fmt.Fprintf(w, str) // отправляем данные на клиентскую сторону
}

func main() {
    http.HandleFunc("/", HomeRouterHandler) // установим роутер
    err := http.ListenAndServe(":9900", nil) // задаем слушать порт
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}