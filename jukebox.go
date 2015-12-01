package main

import (
    "fmt"
    "net/http"
    "github.com/fhs/gompd/mpd"
    "encoding/json"
)

type Song struct {
    ID string `json:"id"`
    Artist string `json:"artist"`
    Title string `json:"title"`
}

const MPDSERVER string = "192.168.0.173:6600"

func playlist(w http.ResponseWriter, r *http.Request) {
    mpdClient, err := mpd.Dial("tcp", MPDSERVER)
    if err != nil {
        fmt.Printf("Error dialing")
    }
    songs, err := mpdClient.PlaylistInfo(-1,-1)
    if err != nil { 
        fmt.Fprintf(w, "Error");
        return
    }
    output := []Song{}
    for _, song := range songs {
        output = append(output, Song{song["Id"], song["Artist"], song["Title"]})
        /*
        for key, value := range song {
            fmt.Fprintf(w, "%s: %s\n", key, value)
        }
        fmt.Fprintf(w, "-----\n")
        */
    }
    js, err := json.Marshal(output)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
    defer mpdClient.Close()
}

func play(w http.ResponseWriter, r *http.Request) {
    mpdClient, err := mpd.Dial("tcp", MPDSERVER)
    if err != nil {
        fmt.Printf("Error dialing")
    }

    err = mpdClient.Play(-1)
    
    if err != nil { 
        fmt.Fprintf(w, "Error play");
        return
    }
    defer mpdClient.Close()
}

func add(w http.ResponseWriter, r *http.Request) {
    url := r.FormValue("url")
    mpdClient, err := mpd.Dial("tcp", MPDSERVER)
    if err != nil {
        fmt.Printf("Error dialing")
        return
    }

    if url != "" {
        err = mpdClient.Add(url)
        if err != nil { 
            fmt.Fprintf(w, "Error adding %s", url);
            return
        }
        fmt.Fprintf(w, "ok")
    }
    
    defer mpdClient.Close()

}

func main() {
    http.HandleFunc("/playlist/add", add)
    http.HandleFunc("/playlist", playlist)
    http.HandleFunc("/play", play)
    http.ListenAndServe(":8080", nil)
}
