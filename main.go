package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func main() {
	// Parse CLI Arguments
	inputFlag := flag.String("i", "", "Input path of image or video to remove watermark")
	outputFlag := flag.String("o", "", "Output path of image or video (optional, defaults to overwrite in-place)")
	flag.Parse()

	// 1. Check if flags are used
	if *inputFlag != "" {
		runCLIMode(*inputFlag, *outputFlag)
		return
	}

	// 2. Fallback to positional arguments: ./tool [input] [output]
	if flag.NArg() >= 1 {
		inputPath := flag.Arg(0)
		outputPath := ""
		if flag.NArg() >= 2 {
			outputPath = flag.Arg(1)
		}
		runCLIMode(inputPath, outputPath)
		return
	}

	// Print usage if no arguments are provided
	flag.Usage()
}

func runCLIMode(inputPath, outputPath string) {
	log.Printf("🚀 Running CLI Watermark Remover on: %s", inputPath)

	info, err := os.Stat(inputPath)
	if os.IsNotExist(err) {
		log.Fatalf("❌ Error: Input file or directory does not exist: %s", inputPath)
	}

	if info.IsDir() {
		// 1. Process directory of files
		files, err := os.ReadDir(inputPath)
		if err != nil {
			log.Fatalf("❌ Failed to read input directory: %v", err)
		}

		// Ensure output directory exists if provided
		if outputPath != "" {
			err = os.MkdirAll(outputPath, 0755)
			if err != nil {
				log.Fatalf("❌ Failed to create output directory: %v", err)
			}
		}

		processedCount := 0
		startTime := time.Now()

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			filename := file.Name()
			ext := strings.ToLower(filepath.Ext(filename))
			var fileType string
			if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
				fileType = "image"
			} else if ext == ".mp4" {
				fileType = "video"
			} else {
				continue // Skip unsupported files
			}

			inputFile := filepath.Join(inputPath, filename)
			var outputFile string
			if outputPath != "" {
				outputFile = filepath.Join(outputPath, "cleaned_"+filename)
			} else {
				// Default behavior: create a cleaned version in the same directory and delete the original
				outputFile = filepath.Join(inputPath, "cleaned_"+filename)
			}

			absInput, err := filepath.Abs(inputFile)
			if err != nil {
				log.Printf("⚠️ Failed to parse absolute path for %s: %v", filename, err)
				continue
			}

			absOutput, err := filepath.Abs(outputFile)
			if err != nil {
				log.Printf("⚠️ Failed to parse absolute output path for %s: %v", filename, err)
				continue
			}

			if absOutput != absInput {
				err = copyFile(absInput, absOutput)
				if err != nil {
					log.Printf("⚠️ Failed to copy file %s: %v", filename, err)
					continue
				}
			}

			fileStartTime := time.Now()
			log.Printf("🧹 Removing watermark from %s (%s)...", filename, fileType)
			err = RemoveWatermark(absOutput, fileType)
			if err != nil {
				log.Printf("❌ Failed to process %s: %v", filename, err)
				if absOutput != absInput {
					os.Remove(absOutput) // clean up failed output
				}
				continue
			}

			log.Printf("✨ Cleaned %s in %.2fs", filename, time.Since(fileStartTime).Seconds())
			processedCount++

			// Delete original input file if a separate output path was successfully written
			if absOutput != absInput {
				if err := os.Remove(absInput); err != nil {
					log.Printf("⚠️ Warning: Failed to delete original file %s: %v", filename, err)
				} else {
					log.Printf("🗑️ Original file deleted: %s", absInput)
				}
			}
		}

		log.Printf("🎉 Finished batch processing! Cleaned %d files in %.2fs", processedCount, time.Since(startTime).Seconds())
		return
	}

	// 2. Single file processing
	ext := strings.ToLower(filepath.Ext(inputPath))
	var fileType string
	if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
		fileType = "image"
	} else if ext == ".mp4" {
		fileType = "video"
	} else {
		log.Fatalf("❌ Error: Unsupported file format '%s'. Must be .png, .jpg, .jpeg, or .mp4", ext)
	}

	// Determine output path
	finalOutPath := outputPath
	if finalOutPath == "" {
		// Default behavior: prefix "cleaned_" at the start of the filename and delete the original
		dir := filepath.Dir(inputPath)
		base := filepath.Base(inputPath)
		finalOutPath = filepath.Join(dir, "cleaned_"+base)
	}

	// Setup absolute path
	absInput, err := filepath.Abs(inputPath)
	if err != nil {
		log.Fatalf("❌ Failed to parse absolute path: %v", err)
	}

	// Setup working path (if writing to a different path, we copy input to target first, then clean)
	if finalOutPath != inputPath {
		err := copyFile(absInput, finalOutPath)
		if err != nil {
			log.Fatalf("❌ Failed to create output file: %v", err)
		}
	}

	absOutput, err := filepath.Abs(finalOutPath)
	if err != nil {
		log.Fatalf("❌ Failed to parse output path: %v", err)
	}

	// Process watermark removal in-place on target path
	startTime := time.Now()
	err = RemoveWatermark(absOutput, fileType)
	if err != nil {
		log.Fatalf("❌ Failed to remove watermark: %v", err)
	}

	// Delete the original input file if a separate output path was successfully written
	if finalOutPath != inputPath {
		if err := os.Remove(absInput); err != nil {
			log.Printf("⚠️ Warning: Failed to delete original input file: %v", err)
		} else {
			log.Printf("🗑️ Original input file deleted: %s", absInput)
		}
	}

	log.Printf("✨ Success! Cleaned file: %s (Time: %.2fs)", absOutput, time.Since(startTime).Seconds())
}

// copyFile is a helper to copy file content from src to dst
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Sync()
}

// Constants for mathematical reverse alpha blending
const (
	alphaThreshold  = 0.002
	maxAlpha        = 0.99
	logoValue       = 255.0
	videoAlphaScale = 0.6
)

// Embedded base64-encoded PNG assets for the 48px and 96px watermarks
const (
	bg48B64 = "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"
	bg96B64 = "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"
)

// Global cache for decoded alpha maps
var (
	alphaCache48 []float32
	alphaCache96 []float32
)

func init() {
	// Pre-decode and cache the watermark grids
	var err error
	alphaCache48, err = loadAlphaMap(48)
	if err != nil {
		log.Printf("⚠️ Failed to load 48px alpha map: %v", err)
	}
	alphaCache96, err = loadAlphaMap(96)
	if err != nil {
		log.Printf("⚠️ Failed to load 96px alpha map: %v", err)
	}
}

// loadAlphaMap decodes base64 PNGs and extracts the maximum RGB channel as alpha mapping
func loadAlphaMap(size int) ([]float32, error) {
	b64Str := bg48B64
	if size == 96 {
		b64Str = bg96B64
	}

	data, err := base64.StdEncoding.DecodeString(b64Str)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	// Resize using simple nearest-neighbor if dimensions mismatch
	alphaMap := make([]float32, size*size)
	for y := 0; y < size; y++ {
		srcY := int(float64(y) * float64(h) / float64(size))
		for x := 0; x < size; x++ {
			srcX := int(float64(x) * float64(w) / float64(size))
			r, g, b, _ := img.At(bounds.Min.X+srcX, bounds.Min.Y+srcY).RGBA()

			// Max of RGB, scale [0, 65535] down to [0.0, 1.0]
			maxVal := r
			if g > maxVal {
				maxVal = g
			}
			if b > maxVal {
				maxVal = b
			}

			alphaMap[y*size+x] = float32(maxVal) / 65535.0
		}
	}

	return alphaMap, nil
}

// WatermarkConfig holds dimensions and offset coordinates
type WatermarkConfig struct {
	LogoSize int
	X        int
	Y        int
}

// DetectWatermarkConfig returns position configurations based on dimensions
func DetectWatermarkConfig(width, height int, isVideo bool) WatermarkConfig {
	var logoSize, marginRight, marginBottom int

	if isVideo {
		shortDim := width
		if height < shortDim {
			shortDim = height
		}

		if shortDim >= 1080 {
			logoSize = 96
			marginRight = 64
			marginBottom = 64
		} else {
			logoSize = 48
			marginRight = 72
			marginBottom = 72
		}
	} else {
		// Portrait standard matching
		if width == 720 && height == 1280 {
			logoSize = 48
			marginRight = 72
			marginBottom = 72
		} else {
			shortDim := width
			if height < shortDim {
				shortDim = height
			}

			if shortDim > 800 {
				logoSize = 96
				marginRight = 64
				marginBottom = 64
			} else {
				logoSize = 48
				marginRight = 32
				marginBottom = 32
			}
		}
	}

	x := width - marginRight - logoSize
	y := height - marginBottom - logoSize

	return WatermarkConfig{
		LogoSize: logoSize,
		X:        x,
		Y:        y,
	}
}

// RemoveWatermark native hub. Dispatches based on file extensions.
func RemoveWatermark(savePath string, fileType string) error {
	ext := strings.ToLower(filepath.Ext(savePath))
	if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
		return removeWatermarkFromImage(savePath)
	} else if ext == ".mp4" {
		return removeWatermarkFromVideo(savePath)
	}
	return fmt.Errorf("unsupported watermark file type: %s", ext)
}

// FindWatermarkOffset locates the watermark coordinates in an image using cross-correlation template matching
func FindWatermarkOffset(img image.Image, alphaMap []float32, logoSize int, standardX, standardY int, radius int) (int, int, float64) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	bestX, bestY := standardX, standardY

	if width < 200 || height < 200 {
		return bestX, bestY, 0.0
	}

	// Search range of +/- radius pixels around standard guess
	minX := standardX - radius
	maxX := standardX + radius
	minY := standardY - radius
	maxY := standardY + radius

	if minX < 0 {
		minX = 0
	}
	if maxX > width-logoSize {
		maxX = width - logoSize
	}
	if minY < 0 {
		minY = 0
	}
	if maxY > height-logoSize {
		maxY = height - logoSize
	}

	sumAlpha := 0.0
	for _, a := range alphaMap {
		sumAlpha += float64(a)
	}
	meanAlpha := sumAlpha / float64(logoSize*logoSize)

	varAlpha := 0.0
	for _, a := range alphaMap {
		diff := float64(a) - meanAlpha
		varAlpha += diff * diff
	}

	bestScore := -1.0
	patchGray := make([]float64, logoSize*logoSize)

	for cy := minY; cy <= maxY; cy++ {
		for cx := minX; cx <= maxX; cx++ {
			sumPixel := 0.0
			for dy := 0; dy < logoSize; dy++ {
				for dx := 0; dx < logoSize; dx++ {
					r, g, b, _ := img.At(bounds.Min.X+cx+dx, bounds.Min.Y+cy+dy).RGBA()
					// Standard grayscale formula (OpenCV COLOR_BGR2GRAY equivalent)
					gray := (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 257.0
					patchGray[dy*logoSize+dx] = gray
					sumPixel += gray
				}
			}
			meanPixel := sumPixel / float64(logoSize*logoSize)

			varPixel := 0.0
			cov := 0.0
			for i := 0; i < logoSize*logoSize; i++ {
				diffPixel := patchGray[i] - meanPixel
				diffAlpha := float64(alphaMap[i]) - meanAlpha
				varPixel += diffPixel * diffPixel
				cov += diffPixel * diffAlpha
			}

			score := 0.0
			if varPixel > 0.0001 && varAlpha > 0.0001 {
				score = cov / math.Sqrt(varPixel*varAlpha)
			}

			if score > bestScore {
				bestScore = score
				bestX, bestY = cx, cy
			}
		}
	}
	return bestX, bestY, bestScore
}

// FindWatermarkOffsetVideo locates the watermark coordinates in raw Y channel using cross-correlation template matching
func FindWatermarkOffsetVideo(yChannel []byte, width, height int, alphaMap []float32, logoSize int, standardX, standardY int) (int, int, float64) {
	bestX, bestY := standardX, standardY

	minX := standardX - 80
	maxX := standardX + 80
	minY := standardY - 80
	maxY := standardY + 80

	if minX < 0 {
		minX = 0
	}
	if maxX > width-logoSize {
		maxX = width - logoSize
	}
	if minY < 0 {
		minY = 0
	}
	if maxY > height-logoSize {
		maxY = height - logoSize
	}

	sumAlpha := 0.0
	for _, a := range alphaMap {
		sumAlpha += float64(a)
	}
	meanAlpha := sumAlpha / float64(logoSize*logoSize)

	varAlpha := 0.0
	for _, a := range alphaMap {
		diff := float64(a) - meanAlpha
		varAlpha += diff * diff
	}

	bestScore := -1.0
	patchGray := make([]float64, logoSize*logoSize)

	for cy := minY; cy <= maxY; cy++ {
		for cx := minX; cx <= maxX; cx++ {
			sumPixel := 0.0
			for dy := 0; dy < logoSize; dy++ {
				for dx := 0; dx < logoSize; dx++ {
					yVal := float64(yChannel[(cy+dy)*width+(cx+dx)])
					patchGray[dy*logoSize+dx] = yVal
					sumPixel += yVal
				}
			}
			meanPixel := sumPixel / float64(logoSize*logoSize)

			varPixel := 0.0
			cov := 0.0
			for i := 0; i < logoSize*logoSize; i++ {
				diffPixel := patchGray[i] - meanPixel
				diffAlpha := float64(alphaMap[i]) - meanAlpha
				varPixel += diffPixel * diffPixel
				cov += diffPixel * diffAlpha
			}

			score := 0.0
			if varPixel > 0.0001 && varAlpha > 0.0001 {
				score = cov / math.Sqrt(varPixel*varAlpha)
			}

			if score > bestScore {
				bestScore = score
				bestX, bestY = cx, cy
			}
		}
	}
	return bestX, bestY, bestScore
}

func v2SmallConfigFromDims(w, h int) (margin, logoSize int) {
	longSide := w
	if h > longSide {
		longSide = h
	}
	shortSide := w
	if h < shortSide {
		shortSide = h
	}

	var sourceLongDim float64
	if shortSide >= 566 {
		sourceLongDim = 2752.0
	} else if shortSide >= 550 {
		sourceLongDim = 2816.0
	} else {
		sourceLongDim = 2848.0
	}

	scale := float64(longSide) / sourceLongDim
	margin = int(math.Round(192.0 * scale))
	return margin, 36
}

func getAlphaMapForSize(size int) ([]float32, error) {
	if size == 48 {
		return alphaCache48, nil
	}
	if size == 96 {
		return alphaCache96, nil
	}
	return loadAlphaMap(size)
}

// removeWatermarkFromImage performs native Go reverse-alpha editing on static images
func removeWatermarkFromImage(imagePath string) error {
	file, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return err
	}
	file.Close() // Close early for in-place write

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Define Candidates
	type Candidate struct {
		logoSize int
		x        int
		y        int
		variant  string // "V1" or "V2"
	}

	var candidates []Candidate
	isLarge := width > 1024 && height > 1024

	if isLarge {
		// V2 Large
		candidates = append(candidates, Candidate{
			logoSize: 96,
			x:        width - 192 - 96,
			y:        height - 192 - 96,
			variant:  "V2",
		})
		// V1 Large
		candidates = append(candidates, Candidate{
			logoSize: 96,
			x:        width - 64 - 96,
			y:        height - 64 - 96,
			variant:  "V1",
		})
	} else {
		// V1 Small
		candidates = append(candidates, Candidate{
			logoSize: 48,
			x:        width - 32 - 48,
			y:        height - 32 - 48,
			variant:  "V1",
		})
		// V2 Small
		v2Margin, v2Logo := v2SmallConfigFromDims(width, height)
		candidates = append(candidates, Candidate{
			logoSize: v2Logo,
			x:        width - v2Margin - v2Logo,
			y:        height - v2Margin - v2Logo,
			variant:  "V2",
		})
	}

	var bestCand Candidate
	bestX, bestY := 0, 0
	bestScore := -1.0
	var bestAlphaMap []float32

	for _, cand := range candidates {
		alphaMap, err := getAlphaMapForSize(cand.logoSize)
		if err != nil {
			log.Printf("⚠️ Warning: Failed to load alpha map for size %d: %v", cand.logoSize, err)
			continue
		}

		cx, cy, score := FindWatermarkOffset(img, alphaMap, cand.logoSize, cand.x, cand.y, 4)
		if score > bestScore {
			bestScore = score
			bestX, bestY = cx, cy
			bestCand = cand
			bestAlphaMap = alphaMap
		}
	}

	// If the best score is below detection threshold (e.g. 0.25), we skip processing
	const kDetectionThreshold = 0.25
	if bestScore < kDetectionThreshold {
		log.Printf("⚠️ No watermark detected (best score: %.2f), skipping image: %s", bestScore, imagePath)
		return nil
	}

	log.Printf("🎯 Detected %s watermark (size: %d) at (%d, %d) with score %.4f",
		bestCand.variant, bestCand.logoSize, bestX, bestY, bestScore)

	// Create writable canvas
	canvas := image.NewRGBA(bounds)
	draw.Draw(canvas, bounds, img, bounds.Min, draw.Src)

	// Apply math formula on watermark region
	for dy := 0; dy < bestCand.logoSize; dy++ {
		py := bestY + dy
		if py < bounds.Min.Y || py >= bounds.Max.Y {
			continue
		}

		for dx := 0; dx < bestCand.logoSize; dx++ {
			px := bestX + dx
			if px < bounds.Min.X || px >= bounds.Max.X {
				continue
			}

			alpha := bestAlphaMap[dy*bestCand.logoSize+dx]
			if alpha < alphaThreshold {
				continue
			}

			// Clamp alpha to avoid division by zero
			if alpha > maxAlpha {
				alpha = maxAlpha
			}

			r, g, b, a := canvas.At(px, py).RGBA()

			// Scale down values to standard float32 [0.0, 255.0]
			fR := float64(r) / 257.0
			fG := float64(g) / 257.0
			fB := float64(b) / 257.0

			// Solve original = (watermarked - alpha * 255.0) / (1 - alpha)
			oneMinusAlpha := 1.0 - float64(alpha)
			newR := (fR - float64(alpha)*logoValue) / oneMinusAlpha
			newG := (fG - float64(alpha)*logoValue) / oneMinusAlpha
			newB := (fB - float64(alpha)*logoValue) / oneMinusAlpha

			// Clamp to [0, 255]
			cR := uint8(math.Min(math.Max(newR, 0), 255))
			cG := uint8(math.Min(math.Max(newG, 0), 255))
			cB := uint8(math.Min(math.Max(newB, 0), 255))
			cA := uint8(a / 257)

			canvas.SetRGBA(px, py, color.RGBA{R: cR, G: cG, B: cB, A: cA})
		}
	}

	// Save back in-place
	outFile, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if format == "png" {
		return png.Encode(outFile, canvas)
	}
	return jpeg.Encode(outFile, canvas, &jpeg.Options{Quality: 95})
}

// VideoStreamInfo is mapped from ffprobe JSON output
type VideoStreamInfo struct {
	Streams []struct {
		Width      int    `json:"width"`
		Height     int    `json:"height"`
		CodecType  string `json:"codec_type"`
		RFrameRate string `json:"r_frame_rate"`
	} `json:"streams"`
}

// removeWatermarkFromVideo runs native FFMPEG byte pipes for zero-Python video watermark removal
func removeWatermarkFromVideo(videoPath string) error {
	absPath, err := filepath.Abs(videoPath)
	if err != nil {
		return err
	}

	// 1. Get dimensions using ffprobe
	probeCmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_streams", absPath)
	probeOut, err := probeCmd.Output()
	if err != nil {
		return fmt.Errorf("ffprobe failed (is ffmpeg installed?): %w", err)
	}

	var info VideoStreamInfo
	if err := json.Unmarshal(probeOut, &info); err != nil {
		return fmt.Errorf("failed to parse ffprobe json: %w", err)
	}

	var width, height int
	var fps string
	for _, stream := range info.Streams {
		if stream.CodecType == "video" {
			width = stream.Width
			height = stream.Height
			fps = stream.RFrameRate
			break
		}
	}

	if width == 0 || height == 0 {
		return fmt.Errorf("could not extract video dimensions for: %s", videoPath)
	}

	if fps == "" {
		fps = "24"
	}

	config := DetectWatermarkConfig(width, height, true)
	alphaMap := alphaCache48
	if config.LogoSize == 96 {
		alphaMap = alphaCache96
	}

	// Setup temp path for safe in-place rewrite
	tempOut := absPath + ".tmp.mp4"
	defer os.Remove(tempOut)

	// Dynamic cross-platform FFmpeg hardware acceleration and encoding settings
	var readArgs []string
	var encoderArgs []string

	if runtime.GOOS == "darwin" {
		// macOS: Use hardware-accelerated VideoToolbox decoding and encoding
		readArgs = []string{"-hwaccel", "videotoolbox", "-i", absPath, "-f", "rawvideo", "-pix_fmt", "yuv420p", "-v", "quiet", "-"}
		encoderArgs = []string{"-c:v", "h264_videotoolbox", "-b:v", "4000k"}
	} else if isEncoderAvailable("h264_nvenc") {
		// Windows/Linux with Nvidia GPU: Use hardware-accelerated NVENC
		readArgs = []string{"-hwaccel", "cuda", "-i", absPath, "-f", "rawvideo", "-pix_fmt", "yuv420p", "-v", "quiet", "-"}
		encoderArgs = []string{"-c:v", "h264_nvenc", "-preset", "p4", "-b:v", "4000k"}
	} else if isEncoderAvailable("h264_qsv") {
		// Windows/Linux with Intel GPU: Use hardware-accelerated QSV
		readArgs = []string{"-hwaccel", "qsv", "-i", absPath, "-f", "rawvideo", "-pix_fmt", "yuv420p", "-v", "quiet", "-"}
		encoderArgs = []string{"-c:v", "h264_qsv", "-b:v", "4000k"}
	} else {
		// Windows/Linux fallback: Use standard CPU-based H.264 decoding and encoding
		readArgs = []string{"-i", absPath, "-f", "rawvideo", "-pix_fmt", "yuv420p", "-v", "quiet", "-"}
		encoderArgs = []string{"-c:v", "libx264", "-preset", "fast", "-crf", "18"}
	}

	// FFMPEG Read command
	readCmd := exec.Command("ffmpeg", readArgs...)
	stdout, err := readCmd.StdoutPipe()
	if err != nil {
		return err
	}

	// FFMPEG Write command
	writeArgs := []string{
		"-y",
		"-f", "rawvideo", "-pix_fmt", "yuv420p",
		"-s", fmt.Sprintf("%dx%d", width, height), "-r", fps,
		"-i", "-",
		"-i", absPath, // Re-read for audio mapping
		"-map", "0:v", "-map", "1:a?",
	}
	writeArgs = append(writeArgs, encoderArgs...)
	writeArgs = append(writeArgs, "-c:a", "copy", "-movflags", "+faststart", "-v", "quiet", tempOut)

	writeCmd := exec.Command("ffmpeg", writeArgs...)

	stdin, err := writeCmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := readCmd.Start(); err != nil {
		return fmt.Errorf("ffmpeg reader start failed: %w", err)
	}
	if err := writeCmd.Start(); err != nil {
		return fmt.Errorf("ffmpeg writer start failed: %w", err)
	}

	ySize := width * height
	uvSize := (width / 2) * (height / 2)
	frameSize := ySize + 2*uvSize // yuv420p size (W * H * 1.5)
	frameBuf := make([]byte, frameSize)

	isOffsetFound := false

	for {
		_, err := io.ReadFull(stdout, frameBuf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			break
		}

		if !isOffsetFound {
			// Check if frame is too dark (e.g., black or fade-in)
			sumY := 0.0
			for i := 0; i < ySize; i++ {
				sumY += float64(frameBuf[i])
			}
			meanY := sumY / float64(ySize)

			if meanY >= 20.0 { // threshold for non-black frame
				bestX, bestY, score := FindWatermarkOffsetVideo(frameBuf[:ySize], width, height, alphaMap, config.LogoSize, config.X, config.Y)
				if score >= 0.20 {
					config.X = bestX
					config.Y = bestY
					isOffsetFound = true
				}
			}
		}

		// Apply YUV watermark reverse blending directly in raw yuv420p buffer if offset is locked
		if isOffsetFound {
			for dy := 0; dy < config.LogoSize; dy++ {
				py := config.Y + dy
				if py < 0 || py >= height {
					continue
				}

				for dx := 0; dx < config.LogoSize; dx++ {
					px := config.X + dx
					if px < 0 || px >= width {
						continue
					}

					alpha := alphaMap[dy*config.LogoSize+dx]
					if alpha < alphaThreshold {
						continue
					}

					// Scale down watermark intensity for videos
					scaledAlpha := alpha * videoAlphaScale
					if scaledAlpha > maxAlpha {
						scaledAlpha = maxAlpha
					}

					oneMinusAlpha := 1.0 - float64(scaledAlpha)

					// 1. Process Y channel (Luminance)
					yOffset := py*width + px
					if yOffset < ySize {
						yVal := float64(frameBuf[yOffset])
						// Solve original = (watermarked - alpha * 235.0) / (1 - alpha)
						newY := (yVal - float64(scaledAlpha)*235.0) / oneMinusAlpha
						frameBuf[yOffset] = byte(math.Min(math.Max(newY, 0), 255))
					}

					// 2. Process U & V channels (Chroma, subsampled 2x2)
					if py%2 == 0 && px%2 == 0 {
						uvRow := py / 2
						uvCol := px / 2
						uvOffset := uvRow*(width/2) + uvCol

						uOffset := ySize + uvOffset
						vOffset := ySize + uvSize + uvOffset

						if uOffset < ySize+uvSize {
							uVal := float64(frameBuf[uOffset])
							newU := (uVal-128.0)/oneMinusAlpha + 128.0
							frameBuf[uOffset] = byte(math.Min(math.Max(newU, 0), 255))
						}

						if vOffset < frameSize {
							vVal := float64(frameBuf[vOffset])
							newV := (vVal-128.0)/oneMinusAlpha + 128.0
							frameBuf[vOffset] = byte(math.Min(math.Max(newV, 0), 255))
						}
					}
				}
			}
		}

		// Write modified frame to encoder stdin pipe
		if _, err := stdin.Write(frameBuf); err != nil {
			break
		}
	}

	stdin.Close()
	writeCmd.Wait()
	readCmd.Wait()

	// Replace original video with clean copy
	if _, err := os.Stat(tempOut); err == nil {
		return os.Rename(tempOut, absPath)
	}

	return fmt.Errorf("failed to compile clean watermark-free video")
}

// isEncoderAvailable queries FFmpeg to check if a specific hardware-accelerated encoder is supported
func isEncoderAvailable(encoder string) bool {
	cmd := exec.Command("ffmpeg", "-encoders")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), encoder)
}
