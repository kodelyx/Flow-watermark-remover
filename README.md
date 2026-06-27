# 🔷 Gemini Watermark Remover

Remove Google Gemini watermarks from videos and images using Go & FFmpeg.

## 📸 Demo

<div align="center">
  <img src="assets/demo_image.png" alt="Demo Image Output" width="48%" style="border-radius: 8px; box-shadow: 0 4px 12px rgba(0,0,0,0.15);" />
  <img src="assets/demo_video.gif" alt="Demo Video Output" width="48%" style="border-radius: 8px; box-shadow: 0 4px 12px rgba(0,0,0,0.15);" />
</div>

## 🚀 Setup & Usage

```bash
# 1. Compile Go binary
go build -o GeminiWatermarkTool-Go main.go

# 2. Run (Auto-renames file to input_cleaned.mp4 and deletes original)
./GeminiWatermarkTool-Go -i input.mp4

# Save output to custom path (deletes original on success)
./GeminiWatermarkTool-Go -i input.mp4 -o custom_output.mp4
```

## 📋 Requirements
* **macOS:** FFmpeg + `VideoToolbox` (GPU acceleration)
* **Windows / Linux:** FFmpeg + Nvidia (`CUDA`) or Intel (`QSV`) GPU acceleration
