package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"ascii-art/files"
	"ascii-art/printing"
)

// --- Constants ---

const (
	bannersPath = "../banners"
)

// --- Models ---
// READY OR NOT
// PageData passes data to the HTML templates
type PageData struct {
	Result    string
	Error     string
	InputText string // Это поле будет содержать "первозданный" текст
}

// --- Application Struct ---

// application stores dependencies, like templates,
// to avoid using global variables.
type application struct {
	indexTmpl *template.Template
	errorTmpl *template.Template
}

// --- Handlers (now methods on *application) ---

// handleIndex handles GET / and 404s
func (app *application) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.renderError(w, http.StatusNotFound, "404 Not Found")
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.renderError(w, http.StatusMethodNotAllowed, "405 Method Not Allowed")
		return
	}
	app.renderTemplate(w, app.indexTmpl, PageData{})
}

// handleAscii handles GET and POST requests for /ascii
func (app *application) handleAscii(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.renderTemplate(w, app.indexTmpl, PageData{})
	case http.MethodPost:
		app.processAsciiPost(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		app.renderError(w, http.StatusMethodNotAllowed, "4Type Not Allowed")
	}
}

// / processAsciiPost contains the logic for POST /ascii
func (app *application) processAsciiPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.renderError(w, http.StatusBadRequest, "400 Bad Request: Invalid form data.")
		return
	}

	// --- <<< ВОТ ГЛАВНОЕ ИЗМЕНЕНИЕ >>> ---
	// 1. Сохраняем "первозданный" текст
	originalText := r.FormValue("text")
	banner := r.FormValue("banner")

	// 2. Создаем "проработанный" текст для генератора
	processedText := strings.ReplaceAll(originalText, "\\n", "\n")
	// --- <<< КОНЕЦ ИЗМЕНЕНИЯ >>> ---

	// Проверки (кроме isASCII и len) можно делать по originalText
	if originalText == "" || banner == "" {
		// --- ИСПРАВЛЕНО ---
		// Раньше было: app.renderTemplate(w, app.indexTmpl, ...)
		app.renderError(w, http.StatusBadRequest, "400 Bad Request: Text and banner style are required.")
		return
	}

	// А вот isASCII и len нужно проверять по "проработанному" тексту,
	// так как он отражает реальное "намерение" пользователя
	if !isASCII(processedText) {
		// --- ИСПРАВЛЕНО ---
		// Раньше было: app.renderTemplate(w, app.indexTmpl, ...)
		app.renderError(w, http.StatusBadRequest, "400 Bad Request: Text contains non-ASCII characters.")
		return
	}

	if len(processedText) > 1000 {
		// --- ИСПРАВЛЕНО ---
		// Раньше было: app.renderTemplate(w, app.indexTmpl, ...)
		app.renderError(w, http.StatusBadRequest, "400 Bad Request: Text exceeds 1000 characters limit.")
		return
	}

	// Generate art using the external packages
	// <<< ИЗМЕНЕНИЕ: Используем "проработанный" текст
	asciiArt, err := generateASCII(processedText, banner)
	if err != nil {
		log.Printf("ASCII generation error: %v", err)
		// --- ИСПРАВЛЕНО ---
		// Раньше было: app.renderTemplate(w, app.indexTmpl, ...)
		app.renderError(w, http.StatusInternalServerError, "500 Internal Server Error: Failed to generate art.")
		return
	}

	// В случае успеха, возвращаем результат И "первозданный" текст
	app.renderTemplate(w, app.indexTmpl, PageData{
		Result:    asciiArt,
		InputText: originalText, // <<< Возвращаем ОРИГИНАЛ
	})
}

// handleFavicon returns 204 No Content to avoid 404 errors
// This handler has no dependencies, so it can remain a plain function
func handleFavicon(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// --- Helper Functions (now methods on *application) ---

// renderTemplate executes HTML templates
func (app *application) renderTemplate(w http.ResponseWriter, tmpl *template.Template, data PageData) {
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	}
}

// renderError renders the custom error page
func (app *application) renderError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	// Use the errorTmpl from our 'app' struct
	if err := app.errorTmpl.Execute(w, PageData{Error: message}); err != nil {
		log.Printf("Error template execution error: %v", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	}
}

// --- Art Generation (unchanged) ---

// isASCII checks for valid printable ASCII
func isASCII(s string) bool {
	for _, c := range s {
		if c == '\n' || c == '\r' {
			continue
		}
		if c < 32 || c > 126 {
			return false
		}
	}
	return true
}

// generateASCII encapsulates the art generation logic.
func generateASCII(text, banner string) (string, error) {
	// Use filepath.Join for correct path assembly
	bannerFile := filepath.Join(bannersPath, banner+".txt")

	// Use the imported function
	bannerMap, err := files.LoadBanner(bannerFile) //
	if err != nil {
		return "", err
	}

	// Use the imported constant
	blankArt := make([]string, printing.ArtHeight) //
	var charArtWidth int
	for _, art := range bannerMap {
		if len(art) > 0 {
			charArtWidth = len(art[0])
			break
		}
	}
	if charArtWidth == 0 {
		log.Println("Could not determine char width for banner", banner)
		return "", err
	}
	for i := 0; i < printing.ArtHeight; i++ { //
		blankArt[i] = strings.Repeat(" ", charArtWidth)
	}

	// strings.Builder is efficient for building strings
	var result strings.Builder
	lines := strings.Split(text, "\n")

	for lineIdx, line := range lines {
		if lineIdx > 0 {
			result.WriteString("\n")
		}
		if line == "" {
			continue
		}

		outputLines := make([]string, printing.ArtHeight) //
		for _, char := range line {
			var charArt []string
			art, ok := bannerMap[char]

			if !ok {
				charArt = blankArt
			} else {
				charArt = make([]string, printing.ArtHeight) //
				for i := 0; i < printing.ArtHeight; i++ {    //
					if i < len(art) {
						charArt[i] = art[i]
					} else {
						charArt[i] = strings.Repeat(" ", charArtWidth)
					}
				}
			}

			for i := 0; i < printing.ArtHeight; i++ { //
				outputLines[i] += charArt[i]
			}
		}
		result.WriteString(strings.Join(outputLines, "\n"))
	}

	return result.String(), nil
}
