# 🔷 Gemini Watermark Remover

### Remove Google Gemini watermarks from videos and images using native Go and FFmpeg pipelines.

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/github/license/kodelyx/Gemini-Watermark-Remover?style=flat-square&color=ff3366)](LICENSE)

> [!IMPORTANT]
> **GPU & Hardware Acceleration Required:**
> * **macOS:** `VideoToolbox` is mandatory. Ensure FFmpeg is built with VideoToolbox support.
> * **Speed:** CPU processing is not supported due to latency.

---

## 📸 Demo

<div align="center">
  <h4>🖼️ Cleaned Image Output</h4>
  <img src="assets/demo_image.png" alt="Demo Image Output" width="90%" style="border-radius: 8px; margin-bottom: 20px; box-shadow: 0 4px 12px rgba(0,0,0,0.15);" />

  <h4>🎥 Cleaned Video Demonstration</h4>
  <img src="assets/demo_video.gif" alt="Demo Video Output" width="90%" style="border-radius: 8px; margin-bottom: 20px; box-shadow: 0 4px 12px rgba(0,0,0,0.15);" />
</div>

---

## ⚡ Compilation & Setup

```bash
# Compile Go utility binary
go build -o GeminiWatermarkTool-Go main.go
```

---

## 🚀 Usage

```bash
# Auto-renames file to input_cleaned.mp4 and deletes the original
./GeminiWatermarkTool-Go -i input.mp4

# Save output to custom path (deletes original on success)
./GeminiWatermarkTool-Go -i input.mp4 -o custom_output.mp4
```

---

## 📄 License
MIT License.
