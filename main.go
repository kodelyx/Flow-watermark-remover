package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//go:embed index.html
var indexHTML string

func main() {
	// Parse CLI Arguments
	inputFlag := flag.String("i", "", "Input path of image or video to remove watermark")
	outputFlag := flag.String("o", "", "Output path of image or video (optional, defaults to overwrite in-place)")
	flag.Parse()

	// 1. CLI Mode
	if *inputFlag != "" {
		runCLIMode(*inputFlag, *outputFlag)
		return
	}

	// 2. HTTP Server Mode (if no input flag is provided)
	runServerMode()
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
				outputFile = inputFile // In-place
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
		finalOutPath = inputPath // In-place update
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

func runServerMode() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(indexHTML))
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok", "service": "watermark-remover"}`))
	})

	http.HandleFunc("/remove-watermark", handleRemoveWatermark)

	fmt.Println("====================================================")
	fmt.Printf("✨ Native Watermark Remover Server Started ✨\n")
	fmt.Printf("Mode: HTTP API + Web Playground\n")
	fmt.Printf("Port: %s\n", port)
	fmt.Printf("URL:  http://localhost:%s\n", port)
	fmt.Println("====================================================")

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleRemoveWatermark(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form-data (Max 50MB)
	err := r.ParseMultipartForm(50 * 1024 * 1024)
	if err != nil {
		http.Error(w, "Failed to parse form: file too large (Max 50MB)", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing 'file' field in request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	var fileType string
	if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
		fileType = "image"
	} else if ext == ".mp4" {
		fileType = "video"
	} else {
		http.Error(w, "Unsupported file format. Must be .png, .jpg, .jpeg, or .mp4", http.StatusBadRequest)
		return
	}

	// Save upload to temporary file
	tempDir := filepath.Join(".", "output", ".temp")
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir) // cleanup temp files immediately after request finishes

	tempFile, err := os.CreateTemp(tempDir, "watermark_*"+ext)
	if err != nil {
		http.Error(w, "Internal server error creating workspace", http.StatusInternalServerError)
		return
	}
	tempPath := tempFile.Name()
	defer os.Remove(tempPath)
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, "Failed to read uploaded file", http.StatusInternalServerError)
		return
	}
	tempFile.Close() // close file handle so RemoveWatermark can process it safely

	// Perform watermark removal in-place on the temp file
	log.Printf("🧹 Removing watermark from uploaded %s (%s)...", header.Filename, fileType)
	err = RemoveWatermark(tempPath, fileType)
	if err != nil {
		log.Printf("❌ Watermark removal failed: %v", err)
		http.Error(w, "Failed to process watermark: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return cleaned file
	cleanedFile, err := os.Open(tempPath)
	if err != nil {
		http.Error(w, "Failed to read processed file", http.StatusInternalServerError)
		return
	}
	defer cleanedFile.Close()

	// Set original content type
	contentType := header.Header.Get("Content-Type")
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"cleaned_%s\"", header.Filename))

	_, err = io.Copy(w, cleanedFile)
	if err != nil {
		log.Printf("❌ Failed to stream processed file to client: %v", err)
	}
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
