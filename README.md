# 🔓 Gemini Watermark Remover

Zero-dependency Go engine to remove Google Gemini watermarks from images and videos. Runs natively on macOS, Windows, and Linux.

---

## 🛠️ Prerequisites

1. **Go:** Make sure Go (1.20+) is installed on your machine.
2. **FFmpeg:** Required for video processing only.
   - **macOS:** `brew install ffmpeg`
   - **Windows:** Download from [ffmpeg.org](https://ffmpeg.org/download.html) and add the binary to your system `PATH`.
   - **Linux:** `sudo apt install ffmpeg`

---

## 🚀 How to Use

### 1. CLI Mode
Process local images or videos directly from the command line.

```bash
# Process a video in-place (overwrites original)
go run . -i input_video.mp4

# Process an image and save it to a new location
go run . -i input_image.png -o cleaned_image.png
```

### 2. Web UI Mode
Start the local web playground to upload and process files in a browser:

```bash
# Start the local server
go run .
```
Open **`http://localhost:8000`** in your browser.
