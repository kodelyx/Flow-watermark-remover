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
    <title>Watermark Remover</title>
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
            max-width: 620px;
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

        /* File List Container */
        .file-list {
            display: flex;
            flex-direction: column;
            gap: 0.75rem;
            margin-bottom: 1.5rem;
            max-height: 280px;
            overflow-y: auto;
            padding-right: 4px;
            text-align: left;
        }

        .file-list::-webkit-scrollbar {
            width: 6px;
        }
        .file-list::-webkit-scrollbar-track {
            background: rgba(255, 255, 255, 0.01);
            border-radius: 3px;
        }
        .file-list::-webkit-scrollbar-thumb {
            background: rgba(255, 255, 255, 0.1);
            border-radius: 3px;
        }
        .file-list::-webkit-scrollbar-thumb:hover {
            background: rgba(255, 255, 255, 0.2);
        }

        .file-item {
            background: rgba(255, 255, 255, 0.04);
            border: 1px solid rgba(255, 255, 255, 0.08);
            border-radius: 12px;
            padding: 0.85rem 1rem;
            display: flex;
            align-items: center;
            justify-content: space-between;
            transition: all 0.2s ease;
        }

        .file-item:hover {
            border-color: rgba(99, 102, 241, 0.3);
            background: rgba(255, 255, 255, 0.06);
        }

        .file-info {
            display: flex;
            align-items: center;
            gap: 0.75rem;
            overflow: hidden;
            flex: 1;
            margin-right: 1rem;
        }

        .file-meta {
            display: flex;
            flex-direction: column;
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
            font-size: 0.75rem;
            color: var(--text-muted);
            margin-top: 0.1rem;
        }

        .status-badge {
            font-size: 0.75rem;
            padding: 0.25rem 0.6rem;
            border-radius: 20px;
            font-weight: 600;
            display: inline-flex;
            align-items: center;
            gap: 0.3rem;
            white-space: nowrap;
        }

        .status-pending {
            background: rgba(255, 255, 255, 0.06);
            color: var(--text-muted);
        }

        .status-processing {
            background: rgba(99, 102, 241, 0.15);
            color: #a5b4fc;
        }

        .status-success {
            background: rgba(34, 197, 94, 0.15);
            color: #86efac;
        }

        .status-error {
            background: rgba(239, 68, 68, 0.15);
            color: #fca5a5;
        }

        .item-actions {
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }

        .item-btn {
            background: none;
            border: none;
            cursor: pointer;
            font-size: 0.9rem;
            display: flex;
            align-items: center;
            justify-content: center;
            border-radius: 6px;
            padding: 0.25rem;
            transition: all 0.2s;
        }

        .remove-item-btn {
            color: rgba(255, 255, 255, 0.4);
        }

        .remove-item-btn:hover {
            color: #ef4444;
            background: rgba(239, 68, 68, 0.1);
        }

        .download-item-btn {
            color: #86efac;
            background: rgba(34, 197, 94, 0.1);
            padding: 0.35rem 0.6rem;
            font-size: 0.75rem;
            font-weight: 600;
            gap: 0.25rem;
            border: 1px solid rgba(34, 197, 94, 0.2);
            border-radius: 8px;
            text-decoration: none;
            display: inline-flex;
            align-items: center;
        }

        .download-item-btn:hover {
            background: rgba(34, 197, 94, 0.2);
            transform: translateY(-1px);
        }

        .inline-spinner {
            width: 12px;
            height: 12px;
            border: 2px solid rgba(255, 255, 255, 0.3);
            border-radius: 50%;
            border-top-color: currentColor;
            animation: spin 1s linear infinite;
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
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 0.5rem;
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
            width: 20px;
            height: 20px;
            border: 3px solid rgba(255, 255, 255, 0.3);
            border-radius: 50%;
            border-top-color: white;
            animation: spin 1s linear infinite;
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
            <p class="subtitle">Instantly remove watermarks from images & videos</p>

            <div id="alertBox" class="alert alert-error"></div>

            <form id="uploadForm">
                <div class="dropzone" id="dropzone">
                    <span class="dropzone-icon">📥</span>
                    <span class="dropzone-text">Drag & drop your files here, or <span style="color:#818cf8; text-decoration: underline;">browse</span></span>
                    <span class="dropzone-subtext">Supports PNG, JPG, JPEG, MP4</span>
                    <input type="file" id="fileInput" name="file" accept=".png,.jpg,.jpeg,.mp4" multiple>
                </div>

                <div class="file-list" id="fileList">
                    <!-- Files will be dynamically inserted here -->
                </div>

                <button type="submit" class="btn" id="submitBtn" disabled>
                    <span id="btnText">Clean Watermark</span>
                    <div class="spinner"></div>
                </button>
            </form>
        </div>
        <div class="footer">
            Premium Edition
        </div>
    </div>

    <script>
        const dropzone = document.getElementById('dropzone');
        const fileInput = document.getElementById('fileInput');
        const fileList = document.getElementById('fileList');
        const submitBtn = document.getElementById('submitBtn');
        const btnText = document.getElementById('btnText');
        const uploadForm = document.getElementById('uploadForm');
        const alertBox = document.getElementById('alertBox');
        
        let selectedFiles = [];

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
                handleFiles(files);
            }
        });

        // Choose file from file picker
        fileInput.addEventListener('change', () => {
            if (fileInput.files.length > 0) {
                handleFiles(fileInput.files);
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

        function handleFiles(files) {
            alertBox.style.display = 'none';
            const allowedExts = ['.png', '.jpg', '.jpeg', '.mp4'];

            for (let i = 0; i < files.length; i++) {
                const file = files[i];
                const name = file.name.toLowerCase();
                const isValid = allowedExts.some(ext => name.endsWith(ext));

                if (!isValid) {
                    showAlert('Some files were skipped. Only PNG, JPG, JPEG, and MP4 are supported.');
                    continue;
                }

                // Check if file is already added
                const alreadyExists = selectedFiles.some(f => f.file.name === file.name && f.file.size === file.size);
                if (alreadyExists) continue;

                selectedFiles.push({
                    id: 'file-' + Date.now() + '-' + Math.random().toString(36).substr(2, 9),
                    file: file,
                    status: 'pending', // pending, processing, success, error
                    progress: 0,
                    downloadUrl: null,
                    error: null
                });
            }

            renderFiles();
            updateSubmitButton();
        }

        function removeFile(id) {
            selectedFiles = selectedFiles.filter(f => f.id !== id);
            renderFiles();
            updateSubmitButton();
        }

        function updateSubmitButton() {
            const pendingCount = selectedFiles.filter(f => f.status === 'pending').length;
            const processingCount = selectedFiles.filter(f => f.status === 'processing').length;
            
            if (processingCount > 0) {
                submitBtn.disabled = true;
                submitBtn.classList.add('loading');
                btnText.textContent = 'Processing (' + processingCount + '/' + selectedFiles.length + ')';
            } else if (pendingCount > 0) {
                submitBtn.disabled = false;
                submitBtn.classList.remove('loading');
                btnText.textContent = 'Clean ' + pendingCount + ' File' + (pendingCount > 1 ? 's' : '');
            } else {
                submitBtn.disabled = true;
                submitBtn.classList.remove('loading');
                btnText.textContent = 'Clean Watermark';
            }
        }

        function renderFiles() {
            if (selectedFiles.length === 0) {
                fileList.innerHTML = '';
                dropzone.style.display = 'block';
                return;
            }

            dropzone.style.display = 'none';

            let html = '<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.5rem;">' +
                       '<span style="font-size: 0.9rem; font-weight: 600; color: var(--text-muted);">' + selectedFiles.length + ' file' + (selectedFiles.length > 1 ? 's' : '') + ' selected</span>' +
                       '<button type="button" onclick="document.getElementById(\'fileInput\').click()" style="background: none; border: none; color: #818cf8; font-size: 0.85rem; font-weight: 600; cursor: pointer; text-decoration: underline;">+ Add More</button>' +
                       '</div>';

            selectedFiles.forEach(item => {
                const isVideo = item.file.name.toLowerCase().endsWith('.mp4');
                const fileIcon = isVideo ? '🎥' : '🖼️';
                
                let statusBadgeHTML = '';
                let actionHTML = '';

                if (item.status === 'pending') {
                    statusBadgeHTML = '<span class="status-badge status-pending">Pending</span>';
                    actionHTML = '<button type="button" class="item-btn remove-item-btn" onclick="removeFile(\'' + item.id + '\')" title="Remove">✕</button>';
                } else if (item.status === 'processing') {
                    statusBadgeHTML = '<span class="status-badge status-processing"><span class="inline-spinner"></span> Cleaning</span>';
                    actionHTML = '';
                } else if (item.status === 'success') {
                    statusBadgeHTML = '<span class="status-badge status-success">✓ Cleaned</span>';
                    const origName = item.file.name;
                    const dotIdx = origName.lastIndexOf('.');
                    const cleanName = origName.substring(0, dotIdx) + '_clean' + origName.substring(dotIdx);
                    
                    actionHTML = '<a href="' + item.downloadUrl + '" download="' + cleanName + '" class="download-item-btn">📥 Download</a>' +
                                 '<button type="button" class="item-btn remove-item-btn" onclick="removeFile(\'' + item.id + '\')" title="Clear">✕</button>';
                } else if (item.status === 'error') {
                    statusBadgeHTML = '<span class="status-badge status-error" title="' + (item.error || 'Error') + '">⚠ Failed</span>';
                    actionHTML = '<button type="button" class="item-btn remove-item-btn" onclick="removeFile(\'' + item.id + '\')" title="Remove">✕</button>';
                }

                html += '<div class="file-item" id="' + item.id + '">' +
                        '<div class="file-info">' +
                        '<span style="font-size: 1.5rem;">' + fileIcon + '</span>' +
                        '<div class="file-meta">' +
                        '<span class="file-name" title="' + item.file.name + '">' + item.file.name + '</span>' +
                        '<span class="file-size">' + formatBytes(item.file.size) + '</span>' +
                        '</div>' +
                        '</div>' +
                        '<div class="item-actions">' +
                        statusBadgeHTML +
                        actionHTML +
                        '</div>' +
                        '</div>';
            });

            fileList.innerHTML = html;
        }

        function showAlert(msg) {
            alertBox.textContent = msg;
            alertBox.style.display = 'block';
        }

        // Form Submit
        uploadForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            
            const pendingItems = selectedFiles.filter(item => item.status === 'pending');
            if (pendingItems.length === 0) return;

            pendingItems.forEach(item => {
                item.status = 'processing';
            });
            renderFiles();
            updateSubmitButton();

            alertBox.style.display = 'none';

            const uploadPromises = pendingItems.map(async (item) => {
                const formData = new FormData();
                formData.append('file', item.file);

                try {
                    const response = await fetch('/remove-watermark', {
                        method: 'POST',
                        body: formData
                    });

                    if (!response.ok) {
                        const errText = await response.text();
                        throw new Error(errText || 'Failed to remove watermark');
                    }

                    const blob = await response.blob();
                    item.downloadUrl = window.URL.createObjectURL(blob);
                    item.status = 'success';
                    
                    const a = document.createElement('a');
                    a.href = item.downloadUrl;
                    const origName = item.file.name;
                    const dotIdx = origName.lastIndexOf('.');
                    const cleanName = origName.substring(0, dotIdx) + '_clean' + origName.substring(dotIdx);
                    a.download = cleanName;
                    document.body.appendChild(a);
                    a.click();
                    document.body.removeChild(a);

                } catch (err) {
                    item.status = 'error';
                    item.error = err.message;
                    console.error('Error processing file:', item.file.name, err);
                }
            });

            await Promise.all(uploadPromises);
            
            renderFiles();
            updateSubmitButton();

            const failedCount = selectedFiles.filter(item => item.status === 'error').length;
            if (failedCount > 0) {
                showAlert('Finished processing. ' + failedCount + ' file' + (failedCount > 1 ? 's' : '') + ' failed to clean.');
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
