package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	ico "github.com/Kodeworks/golang-image-ico"
	"github.com/disintegration/imaging"
)

func main() {
	http.HandleFunc("/generate", cors(generateHandler))
	log.Println("favicon backend running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// CORS middleware
func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

// POST /generate
func generateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "could not parse multipart form: "+err.Error(), 400)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "image is required: "+err.Error(), 400)
		return
	}
	defer file.Close()

	src, _, err := image.Decode(file)
	if err != nil {
		http.Error(w, "failed to decode image: "+err.Error(), 400)
		return
	}

	src = centerCropSquare(src)

	name := r.FormValue("name")
	if name == "" {
		name = "My App"
	}
	shortName := r.FormValue("short_name")
	if shortName == "" {
		shortName = "App"
	}
	themeColor := r.FormValue("theme_color")
	if themeColor == "" {
		themeColor = "#ffffff"
	}
	backgroundColor := r.FormValue("background_color")
	if backgroundColor == "" {
		backgroundColor = "#ffffff"
	}
	tileColor := r.FormValue("tile_color")
	if tileColor == "" {
		tileColor = themeColor
	}

	// Icon sizes and names
	type iconSpec struct {
		Size int
		Name string
	}
	iconSpecs := []iconSpec{
		{16, "favicon-16x16.png"},
		{32, "favicon-32x32.png"},
		{48, "favicon-48x48.png"},
		{57, "apple-icon-57x57.png"},
		{60, "apple-icon-60x60.png"},
		{72, "apple-icon-72x72.png"},
		{76, "apple-icon-76x76.png"},
		{96, "android-icon-96x96.png"},
		{114, "apple-icon-114x114.png"},
		{120, "apple-icon-120x120.png"},
		{128, "android-icon-128x128.png"},
		{144, "android-icon-144x144.png"},
		{152, "apple-icon-152x152.png"},
		{167, "apple-icon-167x167.png"},
		{180, "apple-icon-180x180.png"},
		{192, "android-icon-192x192.png"},
		{256, "android-icon-256x256.png"},
		{384, "android-icon-384x384.png"},
		{512, "android-icon-512x512.png"},
		{150, "ms-icon-150x150.png"},
		{32, "ms-icon-32x32.png"},
		{70, "ms-icon-70x70.png"},
		{310, "ms-icon-310x310.png"},
	}

	files := map[string][]byte{}

	// PNG icons
	for _, spec := range iconSpecs {
		im := imaging.Resize(src, spec.Size, spec.Size, imaging.Lanczos)
		buf := &bytes.Buffer{}
		if err := png.Encode(buf, im); err != nil {
			http.Error(w, "png encode error: "+err.Error(), 500)
			return
		}
		files[spec.Name] = buf.Bytes()
	}

	// Apple touch icon (standard)
	if b, ok := files["apple-icon-180x180.png"]; ok {
		files["apple-touch-icon.png"] = b
	}

	// Multi-size favicon.ico
	icoImg := imaging.Resize(src, 128, 128, imaging.Lanczos)
	icoBuf := &bytes.Buffer{}
	if err := ico.Encode(icoBuf, icoImg); err != nil {
		http.Error(w, "ico encode error: "+err.Error(), 500)
		return
	}
	files["favicon.ico"] = icoBuf.Bytes()

	// Manifest icons
	type iconEntry struct {
		Src     string `json:"src"`
		Sizes   string `json:"sizes"`
		Type    string `json:"type"`
		Purpose string `json:"purpose,omitempty"`
	}
	var icons []iconEntry
	for _, spec := range iconSpecs {
		if strings.HasSuffix(spec.Name, ".png") {
			icons = append(icons, iconEntry{
				Src:   spec.Name,
				Sizes: fmt.Sprintf("%dx%d", spec.Size, spec.Size),
				Type:  "image/png",
			})
		}
	}
	icons = append(icons, iconEntry{
		Src:     "android-icon-512x512.png",
		Sizes:   "512x512",
		Type:    "image/png",
		Purpose: "any maskable",
	})

	manifest := map[string]interface{}{
		"name":             name,
		"short_name":       shortName,
		"start_url":        "/",
		"display":          "standalone",
		"background_color": backgroundColor,
		"theme_color":      themeColor,
		"icons":            icons,
	}
	manBytes, _ := json.MarshalIndent(manifest, "", "  ")
	files["site.webmanifest"] = manBytes
	files["manifest.json"] = manBytes

	// browserconfig.xml for Windows tiles
	bc := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<browserconfig>
  <msapplication>
    <tile>
      <square70x70logo src="ms-icon-70x70.png"/>
      <square150x150logo src="ms-icon-150x150.png"/>
      <square310x310logo src="ms-icon-310x310.png"/>
      <TileColor>%s</TileColor>
    </tile>
  </msapplication>
</browserconfig>`, tileColor)
	files["browserconfig.xml"] = []byte(bc)

	// README
	files["README.txt"] = []byte(buildReadme(name, themeColor))

	// Create zip
	zipBuf := &bytes.Buffer{}
	zw := zip.NewWriter(zipBuf)
	for path, data := range files {
		fh, _ := zw.Create(path)
		if _, err := fh.Write(data); err != nil {
			zw.Close()
			http.Error(w, "zip write error: "+err.Error(), 500)
			return
		}
	}
	if err := zw.Close(); err != nil {
		http.Error(w, "zip finalize error: "+err.Error(), 500)
		return
	}

	ts := strconv.FormatInt(time.Now().Unix(), 10)
	filename := "favicons-" + ts + ".zip"
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	http.ServeContent(w, r, filename, time.Now(), bytes.NewReader(zipBuf.Bytes()))
}

func centerCropSquare(img image.Image) image.Image {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	min := w
	if h < w {
		min = h
	}
	return imaging.CropCenter(img, min, min)
}

func buildReadme(appName, themeColor string) string {
	return fmt.Sprintf("Generated favicons for %s\n\nAll icons are in the root of the ZIP.\nTheme color: %s\n", appName, themeColor)
}
