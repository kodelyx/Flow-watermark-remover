<h1 align="center">🔷 Gemini Watermark Remover</h1>
<p align="center">
  <b>v0.6.2</b> · Remove Gemini / Veo watermarks from videos<br>
  Single binary · Zero dependencies · GPU accelerated
</p>

<p align="center">
  <a href="https://github.com/kodelyx/Gemini-Watermark-Remover/releases/latest"><img src="https://img.shields.io/github/v/release/kodelyx/Gemini-Watermark-Remover?label=Download&color=brightgreen&style=flat-square" /></a>
  <img src="https://img.shields.io/badge/Platform-macOS%20|%20Windows%20|%20Linux-blue?style=flat-square" />
  <img src="https://img.shields.io/badge/GPU-Vulkan%20(MoltenVK)-green?style=flat-square" />
  <img src="https://img.shields.io/badge/License-MIT-yellow?style=flat-square" />
</p>

---

## ✨ Demo

<p align="center">
  <img src="https://raw.githubusercontent.com/allenk/VeoWatermarkRemover/main/artworks/GWT_VEO_Watermark_Removal_Demo.gif" width="600" alt="Before/After Demo" />
</p>

---

## 📋 Overview

| Key       | Value                                            |
|-----------|--------------------------------------------------|
| Tool      | `GeminiWatermarkTool-Video`                      |
| Version   | `0.6.2`                                          |
| Algorithm | Reverse Alpha Blending (math-based, lossless)    |
| GPU       | Vulkan via MoltenVK (macOS) / native (Win/Linux) |
| AI Denoise| FDnCNN neural network (NCNN + Vulkan GPU)        |
| Speed     | ~8s per video (240 frames @ ~29 fps on M1 GPU)   |

### How it works

```
original = (watermarked - alpha × logo_value) / (1 - alpha)
```

No cloud. No AI hallucination. No quality loss. Just math.

---

## 📦 Download

| Platform         | File                          | Notes                    |
|------------------|-------------------------------|--------------------------|
| **macOS**        | `GeminiWatermarkTool-Video`   | Universal (Intel + M1/M2/M3) |
| **Windows x64**  | `GeminiWatermarkTool-Video.exe` | Drag & drop supported  |
| **Linux x64**    | `GeminiWatermarkTool-Video`   | `chmod +x` before running |

👉 [**Download Latest Release**](https://github.com/kodelyx/Gemini-Watermark-Remover/releases/latest)

---

## 🚀 Quick Start

### Step 1 — First-time setup (macOS)

```bash
chmod +x GeminiWatermarkTool-Video
xattr -dr com.apple.quarantine GeminiWatermarkTool-Video
```

### Step 2 — Remove watermark

```bash
# Gemini 3.5 diamond watermark (default)
./GeminiWatermarkTool-Video --veo -i input.mp4 -o clean.mp4

# Older "Veo" text watermark (pre-Gemini 3.5)
./GeminiWatermarkTool-Video --veo --legacy -i input.mp4 -o clean.mp4

# Batch — process all videos in a folder
for f in videos/*.mp4; do
  ./GeminiWatermarkTool-Video --veo -i "$f" -o "clean/$(basename "$f")"
done
```

---

## ⚡ GPU Setup

> Without GPU setup, the tool still works — just slower (CPU fallback).

### macOS (Apple Silicon / Intel)

macOS doesn't ship Vulkan natively. Install MoltenVK:

```bash
# 1. Install via Homebrew
brew install vulkan-loader vulkan-headers molten-vk

# 2. Set environment variables
export VK_ICD_FILENAMES=/opt/homebrew/etc/vulkan/icd.d/MoltenVK_icd.json
export VK_DRIVER_FILES=/opt/homebrew/etc/vulkan/icd.d/MoltenVK_icd.json
export DYLD_LIBRARY_PATH=/opt/homebrew/opt/molten-vk/lib:/opt/homebrew/opt/vulkan-loader/lib

# 3. (Optional) Make permanent — add to ~/.zshrc
echo 'export VK_ICD_FILENAMES=/opt/homebrew/etc/vulkan/icd.d/MoltenVK_icd.json' >> ~/.zshrc
echo 'export VK_DRIVER_FILES=/opt/homebrew/etc/vulkan/icd.d/MoltenVK_icd.json' >> ~/.zshrc
echo 'export DYLD_LIBRARY_PATH=/opt/homebrew/opt/molten-vk/lib:/opt/homebrew/opt/vulkan-loader/lib' >> ~/.zshrc
```

**Verify GPU is active** — run with `--verbose`:
```
NcnnDenoiser: Vulkan GPU #0 - Apple M1
NcnnDenoiser: Initialized (Apple M1)
```

### Windows / Linux

| OS      | Setup                                                                |
|---------|----------------------------------------------------------------------|
| Windows | Works out of the box (any Vulkan-capable GPU)                        |
| Linux   | `sudo apt install mesa-vulkan-drivers` (or distro equivalent)        |

---

## 🔧 CLI Flags

| Flag              | Description                                     |
|-------------------|-------------------------------------------------|
| `--veo`           | Veo video watermark removal mode                |
| `--legacy`        | Old "Veo" text profile (pre-Gemini 3.5)         |
| `--verbose`       | Detailed logs (GPU info, per-frame stats)        |
| `--sigma <N>`     | AI denoise strength (default 50, lower=sharper)  |
| `--variant <V>`   | Force 720p profile: `720p-1` or `720p-2`        |
| `--denoise <M>`   | Cleanup method: `ai`, `ns`, `telea`, `soft`, `off` |
| `-i <file>`       | Input video path                                |
| `-o <file>`       | Output video path                               |

---

## 🛡️ First Run — OS Security

<details>
<summary><b>macOS</b> — "Apple cannot check it for malicious software"</summary>

```bash
xattr -dr com.apple.quarantine GeminiWatermarkTool-Video
chmod +x GeminiWatermarkTool-Video
```
</details>

<details>
<summary><b>Windows</b> — SmartScreen warning</summary>

Click **More info** → **Run anyway**, or in PowerShell:
```powershell
Unblock-File .\GeminiWatermarkTool-Video.exe
```
</details>

---

## 📄 License

MIT — Based on [GeminiWatermarkTool](https://github.com/allenk/GeminiWatermarkTool) by Allen Kuo.
