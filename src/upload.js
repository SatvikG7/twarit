/**
 * File Upload Manager
 * Handles file selection, upload progress, and UI updates
 */
class FileUploader {
    constructor() {
        this.files = [];
        this.isUploading = false;
        this.sessionId = null;
        this.uploadStartTime = null;

        // Cache DOM elements
        this.elements = {};

        this.init();
    }

    init() {
        this.cacheElements();
        this.attachEventListeners();
        this.updateUI();
    }

    cacheElements() {
        const elementIds = [
            "drop-zone",
            "file-input",
            "file-list",
            "upload-btn",
            "progress-container",
            "progress-fill",
            "progress-percentage",
            "upload-speed",
            "file-progress",
            "time-remaining",
            "current-file",
            "status-message",
        ];

        elementIds.forEach((id) => {
            this.elements[id] = document.getElementById(id);
        });

        // Cache common elements
        this.elements.dropZoneIcon =
            this.elements["drop-zone"].querySelector(".drop-zone-icon");
        this.elements.dropZoneText =
            this.elements["drop-zone"].querySelector(".drop-zone-text");
        this.elements.dropZoneSubtext =
            this.elements["drop-zone"].querySelector(".drop-zone-subtext");
    }

    attachEventListeners() {
        // Use event delegation for better performance
        document.addEventListener("click", this.handleClick.bind(this));
        document.addEventListener("change", this.handleChange.bind(this));

        // Drop zone events
        this.elements["drop-zone"].addEventListener(
            "dragover",
            this.handleDragOver.bind(this)
        );
        this.elements["drop-zone"].addEventListener(
            "dragleave",
            this.handleDragLeave.bind(this)
        );
        this.elements["drop-zone"].addEventListener(
            "drop",
            this.handleDrop.bind(this)
        );
    }

    handleClick(event) {
        const { target } = event;

        if (target.id === "drop-zone" && !this.isUploading) {
            this.elements["file-input"].click();
        } else if (target.id === "upload-btn") {
            this.startUpload();
        } else if (target.classList.contains("remove-btn")) {
            const index = parseInt(target.dataset.index);
            this.removeFile(index);
        }
    }

    handleChange(event) {
        if (event.target.id === "file-input") {
            const files = Array.from(event.target.files);
            this.addFiles(files);
        }
    }

    handleDragOver(event) {
        event.preventDefault();
        if (!this.isUploading) {
            this.elements["drop-zone"].classList.add("dragover");
        }
    }

    handleDragLeave() {
        this.elements["drop-zone"].classList.remove("dragover");
    }

    handleDrop(event) {
        event.preventDefault();
        this.elements["drop-zone"].classList.remove("dragover");

        if (!this.isUploading) {
            const files = Array.from(event.dataTransfer.files);
            this.addFiles(files);
        }
    }

    addFiles(newFiles) {
        // Filter out duplicates
        const uniqueFiles = newFiles.filter(
            (newFile) =>
                !this.files.some(
                    (existingFile) =>
                        existingFile.name === newFile.name &&
                        existingFile.size === newFile.size
                )
        );

        this.files = [...this.files, ...uniqueFiles];
        this.updateUI();
    }

    removeFile(index) {
        this.files.splice(index, 1);
        this.updateUI();
    }

    // Batch DOM updates to reduce reflows
    updateUI() {
        // Use requestAnimationFrame to batch DOM updates
        requestAnimationFrame(() => {
            this.updateFileList();
            this.updateUploadButton();
            this.updateDropZone();
        });
    }

    updateFileList() {
        const fragment = document.createDocumentFragment();

        this.files.forEach((file, index) => {
            const fileItem = this.createFileItem(file, index);
            fragment.appendChild(fileItem);
        });

        // Single DOM update
        this.elements["file-list"].replaceChildren(fragment);
    }

    createFileItem(file, index) {
        const fileItem = document.createElement("div");
        fileItem.className = "file-item";

        fileItem.innerHTML = `
            <div class="file-info">
                <div class="file-icon">${this.getFileIcon(file.name)}</div>
                <div class="file-details">
                    <div class="file-name">${this.escapeHtml(file.name)}</div>
                    <div class="file-size">${this.formatFileSize(
                        file.size
                    )}</div>
                </div>
            </div>
            <div class="file-actions">
                ${
                    !this.isUploading
                        ? `<button class="remove-btn" data-index="${index}" type="button" aria-label="Remove ${this.escapeHtml(
                              file.name
                          )}">×</button>`
                        : ""
                }
            </div>
        `;

        return fileItem;
    }

    updateUploadButton() {
        const btn = this.elements["upload-btn"];
        btn.disabled = this.files.length === 0 || this.isUploading;

        const fileCount = this.files.length;
        const fileText = fileCount === 1 ? "File" : "Files";

        btn.textContent = this.isUploading
            ? "⏳ Uploading..."
            : `📤 Upload ${fileCount} ${fileText}`;
    }

    updateDropZone() {
        const { dropZoneIcon, dropZoneText, dropZoneSubtext } = this.elements;
        const dropZone = this.elements["drop-zone"];

        if (this.isUploading) {
            dropZoneIcon.textContent = "⏳";
            dropZoneText.textContent = "Upload in progress...";
            dropZoneSubtext.textContent = "Please wait";
            dropZone.classList.add("uploading");
        } else if (this.files.length > 0) {
            dropZoneIcon.textContent = "📁";
            const fileCount = this.files.length;
            const fileText = fileCount === 1 ? "file" : "files";
            dropZoneText.textContent = `${fileCount} ${fileText} ready to upload`;
            dropZoneSubtext.textContent = "Click to add more files";
            dropZone.classList.remove("uploading");
        } else {
            dropZoneIcon.textContent = "📁";
            dropZoneText.textContent = "Drag & Drop your files here";
            dropZoneSubtext.textContent = "or click to browse";
            dropZone.classList.remove("uploading");
        }
    }

    async startUpload() {
        if (this.files.length === 0 || this.isUploading) return;

        this.isUploading = true;
        this.updateUI();
        this.showStatus("uploading", "Preparing upload...");
        this.elements["progress-container"].classList.add("active");

        try {
            const formData = new FormData();
            this.files.forEach((file) => formData.append("file", file));

            const xhr = new XMLHttpRequest();

            xhr.upload.addEventListener(
                "progress",
                this.handleUploadProgress.bind(this)
            );
            xhr.addEventListener("load", this.handleUploadComplete.bind(this));
            xhr.addEventListener("error", () =>
                this.handleUploadError("Network error occurred")
            );

            xhr.open("POST", "/upload");
            this.uploadStartTime = Date.now();
            xhr.send(formData);
        } catch (error) {
            this.handleUploadError(`Upload failed: ${error.message}`);
        }
    }

    handleUploadProgress(event) {
        if (!event.lengthComputable) return;

        const percentage = (event.loaded / event.total) * 100;
        const speed = this.calculateSpeed(event.loaded);
        const remainingTime = this.calculateRemainingTime(
            event.loaded,
            event.total
        );
        const completedFiles = Math.floor(
            (percentage / 100) * this.files.length
        );
        const currentFile = this.files[completedFiles]?.name || "";

        this.updateProgress({
            percentage,
            speed,
            remainingTime,
            completedFiles,
            totalFiles: this.files.length,
            currentFile,
        });
    }

    handleUploadComplete(event) {
        const xhr = event.target;

        if (xhr.status === 200) {
            try {
                const response = JSON.parse(xhr.responseText);
                this.sessionId = response.sessionId;
                this.onUploadSuccess();
            } catch {
                this.handleUploadError("Invalid server response");
            }
        } else {
            this.handleUploadError(`Upload failed: ${xhr.statusText}`);
        }
    }

    updateProgress({
        percentage,
        speed,
        remainingTime,
        completedFiles,
        totalFiles,
        currentFile,
    }) {
        // Batch DOM updates
        requestAnimationFrame(() => {
            this.elements["progress-fill"].style.width = `${Math.min(
                percentage,
                100
            )}%`;
            this.elements["progress-percentage"].textContent = `${Math.round(
                percentage
            )}%`;
            this.elements["upload-speed"].textContent = this.formatSpeed(speed);
            this.elements[
                "file-progress"
            ].textContent = `${completedFiles} / ${totalFiles}`;
            this.elements["time-remaining"].textContent =
                this.formatTime(remainingTime);

            if (currentFile) {
                this.elements[
                    "current-file"
                ].textContent = `Current: ${currentFile}`;
            }
        });
    }

    onUploadSuccess() {
        this.updateProgress({
            percentage: 100,
            speed: 0,
            remainingTime: 0,
            completedFiles: this.files.length,
            totalFiles: this.files.length,
            currentFile: "",
        });

        const fileText = this.files.length === 1 ? "file" : "files";
        this.showStatus(
            "success",
            `🎉 Successfully uploaded ${this.files.length} ${fileText}!`
        );

        setTimeout(() => this.resetUploader(), 3000);
    }

    handleUploadError(message) {
        this.showStatus("error", `❌ ${message}`);
        this.resetUploader();
    }

    showStatus(type, message) {
        const statusEl = this.elements["status-message"];
        statusEl.className = `status-message ${type}`;
        statusEl.textContent = message;
    }

    resetUploader() {
        this.isUploading = false;
        this.sessionId = null;
        this.uploadStartTime = null;
        this.files = [];

        this.elements["progress-container"].classList.remove("active");
        this.elements["status-message"].className = "status-message";

        this.updateUI();
    }

    // Utility methods
    calculateSpeed(uploadedBytes) {
        if (!this.uploadStartTime) return 0;

        const elapsed = (Date.now() - this.uploadStartTime) / 1000;
        return elapsed > 0 ? uploadedBytes / elapsed : 0;
    }

    calculateRemainingTime(uploadedBytes, totalBytes) {
        const speed = this.calculateSpeed(uploadedBytes);
        return speed > 0 ? (totalBytes - uploadedBytes) / speed : 0;
    }

    getFileIcon(filename) {
        const ext = filename.split(".").pop()?.toLowerCase() || "";
        const iconMap = {
            // Images
            jpg: "🖼️",
            jpeg: "🖼️",
            png: "🖼️",
            gif: "🖼️",
            bmp: "🖼️",
            svg: "🖼️",
            webp: "🖼️",
            // Videos
            mp4: "🎥",
            avi: "🎥",
            mov: "🎥",
            wmv: "🎥",
            flv: "🎥",
            mkv: "🎥",
            webm: "🎥",
            // Audio
            mp3: "🎵",
            wav: "🎵",
            flac: "🎵",
            aac: "🎵",
            ogg: "🎵",
            m4a: "🎵",
            // Documents
            pdf: "📄",
            doc: "📝",
            docx: "📝",
            txt: "📝",
            rtf: "📝",
            xls: "📊",
            xlsx: "📊",
            csv: "📊",
            ppt: "📊",
            pptx: "📊",
            // Archives
            zip: "🗜️",
            rar: "🗜️",
            "7z": "🗜️",
            tar: "🗜️",
            gz: "🗜️",
            // Code
            html: "💻",
            css: "💻",
            js: "💻",
            ts: "💻",
            py: "💻",
            java: "💻",
            cpp: "💻",
            c: "💻",
            go: "💻",
            rs: "💻",
        };

        return iconMap[ext] || "📄";
    }

    formatFileSize(bytes) {
        if (bytes === 0) return "0 Bytes";

        const k = 1024;
        const sizes = ["Bytes", "KB", "MB", "GB", "TB"];
        const i = Math.floor(Math.log(bytes) / Math.log(k));

        return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`;
    }

    formatSpeed(bytesPerSecond) {
        return `${this.formatFileSize(bytesPerSecond)}/s`;
    }

    formatTime(seconds) {
        if (seconds === 0 || !Number.isFinite(seconds)) return "--";

        const hours = Math.floor(seconds / 3600);
        const minutes = Math.floor((seconds % 3600) / 60);
        const secs = Math.floor(seconds % 60);

        if (hours > 0) return `${hours}h ${minutes}m ${secs}s`;
        if (minutes > 0) return `${minutes}m ${secs}s`;
        return `${secs}s`;
    }

    escapeHtml(text) {
        const div = document.createElement("div");
        div.textContent = text;
        return div.innerHTML;
    }
}

// Initialize when DOM is ready
document.addEventListener("DOMContentLoaded", () => {
    window.uploader = new FileUploader();
});
