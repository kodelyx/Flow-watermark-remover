# Guide — Gemini Watermark Remover

## 📋 Overview

| Key       | Value                                            |
|-----------|--------------------------------------------------|
| Tool      | `GeminiWatermarkTool-Video`                      |
| Version   | `0.6.2`                                          |
| Algorithm | Reverse Alpha Blending (math-based, not AI)      |
| GPU       | Vulkan via MoltenVK (macOS) / native (Win/Linux) |
| Denoise   | FDnCNN neural network (NCNN + Vulkan)            |
| Speed     | ~8s per video (240 frames @ ~29 fps on M1 GPU)   |

---

## 🚀 Quick Start

### Step 1 — First-time setup (macOS only)

```bash
chmod +x GeminiWatermarkTool-Video
xattr -dr com.apple.quarantine GeminiWatermarkTool-Video
```

### Step 2 — Remove watermark

```bash
# Single video (Gemini 3.5 diamond watermark)
./GeminiWatermarkTool-Video --veo -i input.mp4 -o clean.mp4

# Older "Veo" text watermark (pre-Gemini 3.5)
./GeminiWatermarkTool-Video --veo --legacy -i input.mp4 -o clean.mp4

# Batch — all mp4s in a folder
for f in videos/*.mp4; do
  ./GeminiWatermarkTool-Video --veo -i "$f" -o "clean/$(basename "$f")"
done
```

---

## ⚡ GPU Setup

### macOS (Apple Silicon / Intel)

macOS doesn't ship Vulkan. Install MoltenVK to enable GPU:

```bash
# 1. Install dependencies
brew install vulkan-loader vulkan-headers molten-vk

# 2. Set environment variables (add to ~/.zshrc for permanent setup)
export VK_ICD_FILENAMES=/opt/homebrew/etc/vulkan/icd.d/MoltenVK_icd.json
export VK_DRIVER_FILES=/opt/homebrew/etc/vulkan/icd.d/MoltenVK_icd.json
export DYLD_LIBRARY_PATH=/opt/homebrew/opt/molten-vk/lib:/opt/homebrew/opt/vulkan-loader/lib

# 3. Verify GPU is detected
./GeminiWatermarkTool-Video --veo --verbose -i test.mp4 -o out.mp4
# Look for: "NcnnDenoiser: Vulkan GPU #0 - Apple M1"
```

### Windows / Linux

| OS      | Setup                                                                |
|---------|----------------------------------------------------------------------|
| Windows | Works out of the box (any Vulkan-capable GPU)                        |
| Linux   | `sudo apt install mesa-vulkan-drivers` (or distro equivalent)        |

> **Note:** Without Vulkan, the tool falls back to CPU — slower but fully functional.

---

## 🔧 Useful Flags

| Flag              | Description                                     |
|-------------------|-------------------------------------------------|
| `--veo`           | Enable Veo video watermark removal mode         |
| `--legacy`        | Use old "Veo" text profile (pre-Gemini 3.5)     |
| `--verbose`       | Detailed logs (GPU info, per-frame stats)        |
| `--sigma <N>`     | AI denoise strength (default 50, lower=sharper)  |
| `--variant <V>`   | Force 720p profile: `720p-1` or `720p-2`        |
| `-i <file>`       | Input video path                                |
| `-o <file>`       | Output video path                               |
