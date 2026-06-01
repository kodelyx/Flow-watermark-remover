# 🔓 Gemini Watermark Remover

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.26+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version" />
  <img src="https://img.shields.io/badge/Platform-macOS%20%7C%20Linux%20%7C%20Windows-blue?style=for-the-badge" alt="Platform Support" />
  <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" alt="License" />
</p>

---

**Gemini Watermark Remover** is a high-performance, zero-dependency Go engine designed to natively remove Google Gemini and Veo AI watermarks from images and videos. It operates completely in-memory using mathematical reverse alpha-blending, restoring original pixels in milliseconds.

## 🌟 Why Gemini Watermark Remover?

Unlike typical media editors or scripting solutions:
* **Zero Dependencies:** Compiled as a single static binary. **No Python 3, OpenCV, or NumPy required.**
* **Mathematical Precision:** It reverses Google's logo blending formula to cleanly reconstruct original pixels instead of simply blurring the area:
  $$\text{original} = \frac{\text{watermarked} - \alpha \times \text{logo}}{1 - \alpha}$$
* **Low Footprint / High Speed:** Employs raw **FFMPEG Byte Pipes** to stream video frames directly in RAM. A 10-second video is cleaned in **under 3 seconds** without exhausting Disk I/O.
* **Dual Interface:** Run it as a command line tool (CLI) for batch scripts, or start the built-in HTTP server to use a premium, responsive Web Playground directly in your browser.

---

## ⚡ Interface Modes

### 1. Web Playground UI (Browser)
Run the executable without any flags to spin up a lightweight, highly responsive, glassmorphic Web Playground:
```bash
./gemini-watermark-remover
```
Open **`http://localhost:8000`** in your browser to drag-and-drop your media and download cleaned files instantly.

### 2. Command Line Interface (CLI)
Pass a file path directly from your terminal:
```bash
# Process file in-place (automatically overwrites original file with cleaned version)
./gemini-watermark-remover -i input_video.mp4

# Save to a new location (original file is automatically cleaned up after output is written)
./gemini-watermark-remover -i input_image.png -o cleaned_image.png
```

---

## 🛠️ Installation & Setup

### Prerequisites
* **FFMPEG** (required for **video** processing only; image cleaning is fully native in Go). FFMPEG must be added to your system's environment variables (`PATH`).

### Build from Source
Ensure you have Go 1.21+ installed on your system:
```bash
# Clone the repository
git clone https://github.com/kodelyx/gemini-watermark-remover.git
cd gemini-watermark-remover

# Compile the binary
go build -o gemini-watermark-remover .
```

---

## 📡 API Integration

You can also host this utility in a cloud container and interface with it programmatically.

### `POST /remove-watermark`
Process and clean any image or video via HTTP POST request.

**Parameters:**
* `file`: Multipart form-data file (supports PNG, JPG, JPEG, MP4).

**Curl Example:**
```bash
curl -X POST -F "file=@my_avatar_video.mp4" http://localhost:8000/remove-watermark > cleaned_avatar_video.mp4
```

---

## 📄 License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
