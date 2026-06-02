# 🔓 Gemini Watermark Remover

Zero-dependency Go engine to remove Google Gemini and Veo watermarks from images and videos.

## 🛠️ Prerequisites
- **FFmpeg** (Required for video processing only) must be installed and added to your `PATH`.

## 🚀 How to Run

### 1. Web UI / API Mode (Port 8000)
Run without any flags to start the server:
```bash
go run .
```
Access the web playground at **`http://localhost:8000`**.

### 2. CLI Mode
Process files directly from terminal:
```bash
# Process in-place (overwrites original)
go run . -i input_video.mp4

# Save to a new location
go run . -i input_image.png -o cleaned_image.png
```
