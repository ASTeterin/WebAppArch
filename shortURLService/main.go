package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
)

const port = 9000
const JSONPathKey = "paths"
const defaultURL = "https://golang.org/"
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
    var urlForRedirect = defaultURL
    fmt.Println(url)
    var paths = ParseJSON(ShortUrlPaths)
    for key, value := range paths {
        if key == url {
            fmt.Println(value)
            urlForRedirect = fmt.Sprintf("%v", value)
            break
        }
    }
    http.Redirect(w, r, urlForRedirect, http.StatusSeeOther)
}

func main() {
    http.HandleFunc("/", RouterHandler)
    err := http.ListenAndServe(":" + strconv.Itoa(port), nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}