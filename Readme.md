
# 🦊 favicon-generator-go

Welcome to **favicon-generator-go**! This project provides a fast and easy way to generate favicon files from images. With both Go and Python implementations, plus a user-friendly frontend, it's perfect for web developers who want to automate favicon creation for their sites.

---

## 📦 Features

- **Multi-language Support:** Go and Python versions provided.
- **Web Interface:** Simple frontend to upload images and get favicons.
- **Automatic Cropping:** Option to center-crop images for optimal favicon appearance.
- **Favicon Package:** Generates standard favicon formats (`.ico`, `png` sizes) and a ready-to-use manifest.
- **Instant Download:** Favicons are packaged as a ZIP for quick download.
- **Open Source:** Easy to extend and customize.

---

## 🛠️ Installation

### 1. Clone the Repository

```sh
git clone https://github.com/yourusername/favicon-generator-go.git
cd favicon-generator-go
```

### 2. Go Version

- Ensure you have Go installed (v1.18+ recommended).
- Install dependencies:

```sh
go get github.com/Kodeworks/golang-image-ico
go get github.com/disintegration/imaging
```

- Run the server:

```sh
go run go-version/main.go
```

### 3. Python Version

- Ensure you have Python 3.8+ and [Pillow](https://pillow.readthedocs.io/) and [gradio](https://gradio.app/) installed.

```sh
pip install pillow gradio
```

- Run the app:

```sh
python python-version/app.py
```

### 4. Frontend

- Open `frontend/index.html` in your browser to use the web interface.

---

## ⚡ Usage

### Web Interface

1. Open the `frontend/index.html` file in your browser.
2. Upload an image (PNG, JPG, etc.).
3. The app will generate favicons in multiple formats and sizes.
4. Download the ZIP package containing:
   - Standard favicon files (`favicon.ico`, PNGs)
   - `manifest.json` for web apps

### API Usage (Go/Python)

Both implementations expose HTTP endpoints to upload your image and receive the zipped favicons. See source code for endpoint details.

---

## 🤝 Contributing

Contributions are welcome! To get started:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/your-feature`).
3. Make your changes.
4. Submit a pull request.

---

## 📄 License

This project is licensed under the **MIT License**. See [LICENSE](LICENSE) for details.

---

## 📚 Project Structure

```
favicon-generator-go/
│
├── frontend/
│   ├── index.html
│   └── manifest.json
│
├── go-version/
│   └── main.go
│
├── python-version/
│   └── app.py
│
└── README.md
```

---

## 💡 Topics

Favicon, Go, Python, Web Development, Image Processing, Automation

---

Happy favicon generating! 🎉

---

## License
This project is licensed under the **MIT** License.

---
🔗 GitHub Repo: https://github.com/Pranesh-2005/favicon-generator-go