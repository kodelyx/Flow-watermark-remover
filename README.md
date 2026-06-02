# 🔓 Gemini Watermark Remover

Zero-dependency, high-performance Go engine to remove Google Gemini watermarks from images and videos. 

It is designed to run natively on both desktop machines (via CLI/Web server) and Android mobile devices (via Go Mobile bindings) with maximum efficiency.

---

## 🚀 Key Features

* **Direct YUV420p Processing:** Video frames are processed natively in `yuv420p` format. This completely bypasses the slow and lossy YUV-to-RGB color-space conversion, resulting in **100% pixel-perfect lossless recovery** (no red tinting or background artifacts).
* **Apple Silicon Hardware Acceleration:** Under macOS, the engine automatically leverages Apple VideoToolbox API (`-hwaccel videotoolbox` and `h264_videotoolbox` encoder) to decode/encode video streams directly on dedicated GPU cores, processing a 10s video in **just ~1.0s**.
* **Cross-Platform Compatibility:** Dynamic OS detection automatically uses hardware acceleration on macOS while falling back to standard CPU-based H.264 codecs (`libx264` with `-preset fast`) on Windows and Linux machines.
* **Low Memory & Lightweight Footprint:** The CLI Go process uses a tiny **~5.9 MB** memory footprint, processing pixels in-place inside memory.

---

## 🛠️ Prerequisites
- **FFmpeg** (Required for video processing only) must be installed and added to your `PATH`.

---

## 🚀 How to Run

### 1. Web UI / API Mode (Port 8000)
Run without any flags to start the server:
```bash
go run .
```
Access the responsive web playground at **`http://localhost:8000`**.

### 2. CLI Mode
Process files directly from terminal:
```bash
# Process in-place (overwrites original)
go run . -i input_video.mp4

# Save to a new location
go run . -i input_image.png -o cleaned_image.png
```
