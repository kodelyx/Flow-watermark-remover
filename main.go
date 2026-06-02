package main

import (
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

// Embed a beautiful, premium responsive UI directly in the server for browser interactions!
const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Native Watermark Remover</title>
    <link href="https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;600;800&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg-color: #0b0f19;
            --card-bg: rgba(255, 255, 255, 0.03);
            --border-color: rgba(255, 255, 255, 0.08);
            --primary-glow: linear-gradient(135deg, #6366f1 0%, #a855f7 50%, #ec4899 100%);
            --text-main: #f3f4f6;
            --text-muted: #9ca3af;
        }

        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }

        body {
            font-family: 'Outfit', -apple-system, sans-serif;
            background-color: var(--bg-color);
            color: var(--text-main);
            min-height: 100vh;
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            overflow-x: hidden;
            position: relative;
            padding: 2rem 1rem;
        }

        /* Ambient Glowing Background Elements */
        body::before {
            content: '';
            position: absolute;
            width: 400px;
            height: 400px;
            background: radial-gradient(circle, rgba(99, 102, 241, 0.15) 0%, rgba(0,0,0,0) 70%);
            top: -100px;
            left: -100px;
            z-index: 0;
            pointer-events: none;
        }
        body::after {
            content: '';
            position: absolute;
            width: 450px;
            height: 450px;
            background: radial-gradient(circle, rgba(236, 72, 153, 0.12) 0%, rgba(0,0,0,0) 70%);
            bottom: -150px;
            right: -100px;
            z-index: 0;
            pointer-events: none;
        }

        .container {
            width: 100%;
            max-width: 580px;
            z-index: 10;
        }

        .card {
            background: var(--card-bg);
            backdrop-filter: blur(20px);
            -webkit-backdrop-filter: blur(20px);
            border: 1px solid var(--border-color);
            border-radius: 24px;
            padding: 2.5rem;
            box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
            text-align: center;
            transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        }

        .logo-container {
            margin-bottom: 1.5rem;
            display: inline-block;
        }

        .logo-icon {
            font-size: 3rem;
            background: var(--primary-glow);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            animation: pulse 2s infinite ease-in-out;
        }

        h1 {
            font-size: 2.2rem;
            font-weight: 800;
            letter-spacing: -0.03em;
            background: var(--primary-glow);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            margin-bottom: 0.5rem;
        }

        .subtitle {
            color: var(--text-muted);
            font-size: 1rem;
            margin-bottom: 2rem;
            font-weight: 300;
        }

        .dropzone {
            border: 2px dashed rgba(255, 255, 255, 0.15);
            border-radius: 16px;
            padding: 3rem 1.5rem;
            cursor: pointer;
            transition: all 0.2s ease;
            position: relative;
            background: rgba(255, 255, 255, 0.01);
            margin-bottom: 1.5rem;
        }

        .dropzone:hover, .dropzone.dragover {
            border-color: #6366f1;
            background: rgba(99, 102, 241, 0.03);
            transform: scale(1.01);
        }

        .dropzone-icon {
            font-size: 2.5rem;
            margin-bottom: 1rem;
            display: block;
            color: #818cf8;
        }

        .dropzone-text {
            font-size: 0.95rem;
            color: var(--text-muted);
            font-weight: 400;
        }

        .dropzone-subtext {
            font-size: 0.8rem;
            color: rgba(255, 255, 255, 0.3);
            margin-top: 0.5rem;
        }

        input[type="file"] {
            display: none;
        }

        .file-selected {
            background: rgba(255, 255, 255, 0.05);
            border: 1px solid rgba(255, 255, 255, 0.15);
            border-radius: 12px;
            padding: 1rem;
            margin-bottom: 1.5rem;
            display: none;
            align-items: center;
            justify-content: space-between;
            text-align: left;
        }

        .file-info {
            display: flex;
            align-items: center;
            gap: 0.75rem;
            overflow: hidden;
        }

        .file-name {
            font-size: 0.9rem;
            font-weight: 600;
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
        }

        .file-size {
            font-size: 0.8rem;
            color: var(--text-muted);
        }

        .remove-btn {
            background: none;
            border: none;
            color: #ef4444;
            cursor: pointer;
            font-size: 1.1rem;
            display: flex;
            align-items: center;
            transition: opacity 0.2s;
        }

        .remove-btn:hover {
            opacity: 0.8;
        }

        .btn {
            width: 100%;
            background: var(--primary-glow);
            color: white;
            border: none;
            padding: 1rem 2rem;
            font-size: 1rem;
            font-weight: 600;
            border-radius: 14px;
            cursor: pointer;
            box-shadow: 0 4px 15px rgba(99, 102, 241, 0.4);
            transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
            position: relative;
            overflow: hidden;
        }

        .btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 8px 25px rgba(99, 102, 241, 0.55);
        }

        .btn:active {
            transform: translateY(0);
        }

        .btn:disabled {
            background: rgba(255, 255, 255, 0.1);
            color: var(--text-muted);
            cursor: not-allowed;
            box-shadow: none;
            transform: none;
        }

        /* Loading Spinner */
        .spinner {
            display: none;
            width: 24px;
            height: 24px;
            border: 3px solid rgba(255, 255, 255, 0.3);
            border-radius: 50%;
            border-top-color: white;
            animation: spin 1s linear infinite;
            position: absolute;
            left: 50%;
            top: 50%;
            margin-left: -12px;
            margin-top: -12px;
        }

        .btn.loading {
            color: transparent;
        }
        .btn.loading .spinner {
            display: block;
        }

        .footer {
            margin-top: 2rem;
            text-align: center;
            font-size: 0.8rem;
            color: var(--text-muted);
            font-weight: 300;
        }

        @keyframes pulse {
            0%, 100% { transform: scale(1); opacity: 0.9; }
            50% { transform: scale(1.05); opacity: 1; }
        }

        @keyframes spin {
            to { transform: rotate(360deg); }
        }

        .alert {
            display: none;
            padding: 1rem;
            border-radius: 12px;
            margin-bottom: 1.5rem;
            font-size: 0.9rem;
            text-align: left;
        }

        .alert-error {
            background: rgba(239, 68, 68, 0.1);
            border: 1px solid rgba(239, 68, 68, 0.2);
            color: #fca5a5;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="card">
            <div class="logo-container">
                <span class="logo-icon">✨</span>
            </div>
            <h1>Watermark Remover</h1>
            <p class="subtitle">Remove Google Gemini watermarks natively from your images & videos</p>

            <div id="alertBox" class="alert alert-error"></div>

            <form id="uploadForm">
                <div class="dropzone" id="dropzone">
                    <span class="dropzone-icon">📥</span>
                    <span class="dropzone-text">Drag & drop your file here, or <span style="color:#818cf8; text-decoration: underline;">browse</span></span>
                    <span class="dropzone-subtext">Supports PNG, JPG, JPEG, MP4</span>
                    <input type="file" id="fileInput" name="file" accept=".png,.jpg,.jpeg,.mp4">
                </div>

                <div class="file-selected" id="fileSelected">
                    <div class="file-info">
                        <span id="fileIcon" style="font-size:1.5rem;">📄</span>
                        <div>
                            <div class="file-name" id="fileName">image.png</div>
                            <div class="file-size" id="fileSize">1.2 MB</div>
                        </div>
                    </div>
                    <button type="button" class="remove-btn" id="removeBtn">✕</button>
                </div>

                <button type="submit" class="btn" id="submitBtn" disabled>
                    <span class="btn-text">Clean Watermark</span>
                    <div class="spinner"></div>
                </button>
            </form>
        </div>
        <div class="footer">
            Native Go & FFMPEG Engine • Premium Edition
        </div>
    </div>

    <script>
        const dropzone = document.getElementById('dropzone');
        const fileInput = document.getElementById('fileInput');
        const fileSelected = document.getElementById('fileSelected');
        const fileName = document.getElementById('fileName');
        const fileSize = document.getElementById('fileSize');
        const fileIcon = document.getElementById('fileIcon');
        const removeBtn = document.getElementById('removeBtn');
        const submitBtn = document.getElementById('submitBtn');
        const uploadForm = document.getElementById('uploadForm');
        const alertBox = document.getElementById('alertBox');
        
        let selectedFile = null;

        // Click zone triggers browse
        dropzone.addEventListener('click', () => fileInput.click());

        // Drag events
        ['dragenter', 'dragover'].forEach(eventName => {
            dropzone.addEventListener(eventName, (e) => {
                e.preventDefault();
                dropzone.classList.add('dragover');
            }, false);
        });

        ['dragleave', 'drop'].forEach(eventName => {
            dropzone.addEventListener(eventName, (e) => {
                e.preventDefault();
                dropzone.classList.remove('dragover');
            }, false);
        });

        // Drop file
        dropzone.addEventListener('drop', (e) => {
            const dt = e.dataTransfer;
            const files = dt.files;
            if (files.length > 0) {
                handleFile(files[0]);
            }
        });

        // Choose file from file picker
        fileInput.addEventListener('change', () => {
            if (fileInput.files.length > 0) {
                handleFile(fileInput.files[0]);
            }
        });

        function formatBytes(bytes, decimals = 2) {
            if (bytes === 0) return '0 Bytes';
            const k = 1024;
            const dm = decimals < 0 ? 0 : decimals;
            const sizes = ['Bytes', 'KB', 'MB', 'GB'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
        }

        function handleFile(file) {
            const name = file.name.toLowerCase();
            const allowedExts = ['.png', '.jpg', '.jpeg', '.mp4'];
            const isValid = allowedExts.some(ext => name.endsWith(ext));

            if (!isValid) {
                showAlert('Invalid file type. Please upload a PNG, JPG, JPEG, or MP4.');
                return;
            }

            alertBox.style.display = 'none';
            selectedFile = file;
            fileName.textContent = file.name;
            fileSize.textContent = formatBytes(file.size);

            if (name.endsWith('.mp4')) {
                fileIcon.textContent = '🎥';
            } else {
                fileIcon.textContent = '🖼️';
            }

            dropzone.style.display = 'none';
            fileSelected.style.display = 'flex';
            submitBtn.disabled = false;
        }

        // Clear chosen file
        removeBtn.addEventListener('click', () => {
            selectedFile = null;
            fileInput.value = '';
            dropzone.style.display = 'block';
            fileSelected.style.display = 'none';
            submitBtn.disabled = true;
        });

        function showAlert(msg) {
            alertBox.textContent = msg;
            alertBox.style.display = 'block';
        }

        // Form Submit
        uploadForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            if (!selectedFile) return;

            submitBtn.disabled = true;
            submitBtn.classList.add('loading');
            alertBox.style.display = 'none';

            const formData = new FormData();
            formData.append('file', selectedFile);

            try {
                const response = await fetch('/remove-watermark', {
                    method: 'POST',
                    body: formData
                });

                if (!response.ok) {
                    const errText = await response.text();
                    throw new Error(errText || 'Failed to remove watermark');
                }

                // Download cleaned file
                const blob = await response.blob();
                const url = window.URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                
                // Keep same base filename but append _clean
                const origName = selectedFile.name;
                const dotIdx = origName.lastIndexOf('.');
                const cleanName = origName.substring(0, dotIdx) + '_clean' + origName.substring(dotIdx);
                
                a.download = cleanName;
                document.body.appendChild(a);
                a.click();
                window.URL.revokeObjectURL(url);
                a.remove();
            } catch (err) {
                showAlert('Error: ' + err.message);
            } finally {
                submitBtn.disabled = false;
                submitBtn.classList.remove('loading');
            }
        });
    </script>
</body>
</html>
`

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

	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		log.Fatalf("❌ Error: Input file does not exist: %s", inputPath)
	}

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
