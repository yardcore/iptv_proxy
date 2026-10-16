package main

import (
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strings"
)

const m3uFolder = "/home/content/strm/live"

func main() {
    http.HandleFunc("/m3u", handleM3U)

    log.Println("IPTV proxy running on :6077, serving folder:", m3uFolder)
    if err := http.ListenAndServe(":6077", nil); err != nil {
        log.Fatal(err)
    }
}

func handleM3U(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/x-mpegURL")

    err := filepath.Walk(m3uFolder, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }

        ext := strings.ToLower(filepath.Ext(path))
        if ext == ".m3u" || ext == ".m3u8" {
            data, readErr := os.ReadFile(path)
            if readErr != nil {
                log.Println("Error reading", path, readErr)
                return nil
            }

            // Write file contents
            _, _ = w.Write(data)
            _, _ = w.Write([]byte("\n"))
        }
        return nil
    })

    if err != nil {
        log.Println("Walk error:", err)
    }
}
