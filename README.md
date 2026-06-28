# 🌊 Flow Omni Video Watermark Remover

A high-performance tool to remove watermarks from **videos** using Go & FFmpeg.

## 📸 Demo

<div align="center">
  <img src="assets/demo_video.gif" alt="Demo Video Output" width="80%" style="border-radius: 8px; box-shadow: 0 4px 12px rgba(0,0,0,0.15);" />
</div>

## 🚀 Setup & Usage

```bash
# 1. Compile Go binary
go build -o Flow-Omni-Watermark-Remover main.go

# 2. Run
./Flow-Omni-Watermark-Remover input.mp4

# Save output to custom path
./Flow-Omni-Watermark-Remover input.mp4 custom_output.mp4
```

## 📋 Requirements
* **macOS:** FFmpeg + `VideoToolbox` (GPU acceleration)
* **Windows / Linux:** FFmpeg + Nvidia (`CUDA`) or Intel (`QSV`) GPU acceleration
