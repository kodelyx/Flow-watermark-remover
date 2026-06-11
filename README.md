# 🔷 Gemini Watermark Remover

Remove Gemini watermarks from videos. Single binary, zero dependencies.

## Demo

### Video
![Video Watermark Removal](assets/demo_video.gif)

### Image
![Image Watermark Removal](assets/demo_image.png)

---

## Usage

```bash
# 1. First time setup (macOS)
chmod +x GeminiWatermarkTool-Video
xattr -dr com.apple.quarantine GeminiWatermarkTool-Video

# 2. Remove watermark from video
./GeminiWatermarkTool-Video --veo -i input.mp4 -o clean.mp4

# 3. Batch process
for f in videos/*.mp4; do
  ./GeminiWatermarkTool-Video --veo -i "$f" -o "clean/$(basename "$f")"
done
```

---

## GPU Setup (macOS)

```bash
brew install vulkan-loader vulkan-headers molten-vk

export VK_ICD_FILENAMES=/opt/homebrew/etc/vulkan/icd.d/MoltenVK_icd.json
export VK_DRIVER_FILES=/opt/homebrew/etc/vulkan/icd.d/MoltenVK_icd.json
export DYLD_LIBRARY_PATH=/opt/homebrew/opt/molten-vk/lib:/opt/homebrew/opt/vulkan-loader/lib
```

> Add these `export` lines to `~/.zshrc` for permanent GPU support.

---

## Download

👉 [**Latest Release**](https://github.com/kodelyx/Gemini-Watermark-Remover/releases/latest)
