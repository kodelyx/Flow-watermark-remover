# 🔷 Gemini Watermark Remover

Remove Gemini watermarks from videos using native Go & FFmpeg pipelines.

> [!IMPORTANT]
> **GPU & VideoToolbox Mandatory:**
> * You **MUST only use GPU** (CPU encoding is not supported due to slow speed).
> * On macOS, **VideoToolbox (GPU/Hardware acceleration)** is mandatory. Make sure your FFmpeg is installed with VideoToolbox support.

## Compilation
```bash
go build -o GeminiWatermarkTool-Go main.go
```

## Usage

### 1. Remove Gemini Watermark
```bash
# Process a video in-place (updates the video directly)
./GeminiWatermarkTool-Go -i input.mp4
```

### 2. Add Brand Watermark (Akash Digital Marketing)
Adds a single centered watermark:
```bash
python3 add_brand_watermark.py -i input.mp4 -o output.mp4
```
