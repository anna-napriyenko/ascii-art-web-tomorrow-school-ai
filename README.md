🖋 **ASCII Art Web**

A simple web application written in Go that converts user-input text into ASCII art using different banner styles.

The application runs on a custom-built web server, handles GET and POST requests, renders HTML templates with CSS styling, and provides clear error handling for various edge cases.

---

🚀 **Features**

### 🌐 Web Interface

* Clean and intuitive UI for text input
* Dropdown menu to select ASCII art style
* Instant rendering of generated ASCII art

### 🎨 ASCII Art Styles

The project supports three banner styles:

* `standard`
* `shadow`
* `thinkertoy`

### ⚙️ Custom Server Architecture

* Uses a custom `http.Server` and `http.ServeMux`
* Avoids global handlers for better control and extensibility
* Handles both GET and POST requests

### ❗ Error Handling

The application gracefully handles errors and displays user-friendly pages:

* **400 Bad Request** — invalid input or non-ASCII characters
* **404 Not Found** — unknown route
* **405 Method Not Allowed** — unsupported HTTP method
* **500 Internal Server Error** — server-side issues (e.g., missing banner files)

### 🧩 Clean Project Structure

* `files` — banner loading and file handling
* `printing` — constants and ASCII definitions
* `app` — web server logic and routing

---

🛠 **Technologies Used**

* Backend: Go
* Frontend: HTML, CSS
* Networking: net/http
* Templates: html/template

---

⚙️ **Requirements**

* **Go:** version 1.25.0 or newer

---

📦 **Installation and Run**

### 🔧 Local Setup

1. Make sure all project files are present and the directory structure is preserved
   (especially `app`, `banners`, `files`, `printing`)
2. Navigate to the application directory:

   ```bash
   cd ascii-art-web-stylize/app
   ```
3. Run the web server:

   ```bash
   go run .
   ```
4. Open in your browser:

   ```
   http://localhost:8080
   ```

---

📋 **How to Use**

### 🌐 Via Browser

1. Open `http://localhost:8080`
2. Enter any text (only ASCII characters are supported)
3. Select one of the styles:

   * `standard`
   * `shadow`
   * `thinkertoy`
4. Click **Generate**
5. The generated ASCII art will appear in the **Result** section below

---

### 🔌 Via cURL (or Postman)

You can send requests directly to the `/ascii` endpoint using `curl` or any API client.

* **Endpoint:** `POST /ascii`
* **Body type:** `x-www-form-urlencoded`
* **Parameters:**

  * `text` (string) — text to convert
  * `banner` (string) — one of: `standard`, `shadow`, `thinkertoy`

**Example cURL request:**

```bash
curl -X POST http://localhost:8080/ascii \
     -d "text=Hello World" \
     -d "banner=shadow"
```

---

📌 **Notes**

* Only ASCII characters are supported
* Banner files must be present in the `banners` directory
* The project is designed for educational and demonstration purposes
