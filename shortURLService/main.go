package main

import (
    "encoding/json"
    "fmt"      // пакет для форматированного ввода вывода
    "log"      // пакет для логирования
    "net/http" // пакет для поддержки HTTP протокола
)

const JSONPathKey = "paths"

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

    var m = pathConfiguration.(map[string]interface{})
    var paths = m[JSONPathKey]
    return paths.(map[string]interface{})
}

func RouterHandler(w http.ResponseWriter, r *http.Request) {

    var url = r.URL.Path
    fmt.Println(url)
    var paths = ParseJSON(ShortUrlPaths)
    for key, value := range paths {
        if key == url {
            fmt.Println(value)
            var urlForRed = fmt.Sprintf("%v", value)
            http.Redirect(w, r, urlForRed, http.StatusSeeOther)
            break
        }

    }

    //var str = fmt.Sprintf("%v", paths)
    //fmt.Println(paths)
    //fmt.Fprintf(w, str) // отправляем данные на клиентскую сторону
}

func main() {
    http.HandleFunc("/", RouterHandler) // установим роутер
    err := http.ListenAndServe(":9900", nil) // задаем слушать порт
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}