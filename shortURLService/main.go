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

const ShortUrlPaths = `
    {
        "paths" : {
            "/go-http": "https://golang.org/pkg/net/http/",
            "/go-gophers" : "https://github.com/shalakhin/gophericons/blob/master/preview.jpg"
        }
    }`


func ParseJSON(JSONData string) map[string]interface{} {
    var pathConfiguration interface{}
    err := json.Unmarshal([]byte(JSONData), &pathConfiguration)
    if err != nil {
        log.Println(err)
    }
    //var str = con.Paths.p["/go-http"]

    var m = pathConfiguration.(map[string]interface{})
    var paths = m["paths"]
    return paths.(map[string]interface{})
}

func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {

    var sss = ParseJSON(ShortUrlPaths)
    for _, value := range sss {
        fmt.Println(value)
    }
    //var str = fmt.Sprintf("%v", paths)
    //fmt.Println(paths)
    //fmt.Fprintf(w, str) // отправляем данные на клиентскую сторону
}

func main() {
    http.HandleFunc("/", HomeRouterHandler) // установим роутер
    err := http.ListenAndServe(":9900", nil) // задаем слушать порт
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}