// src/components/upload/UploadForm.tsx
import { useState, useCallback } from "react";
import { useUploadPhoto } from "@/hooks/usePhotoQueries";
import { UploadFormProps } from "@/types";
import "@/css/Upload.css";

export function UploadForm({ onComplete }: UploadFormProps) {
  const [files, setFiles] = useState<FileList | null>(null);
  const [uploadProgress, setUploadProgress] = useState<{
    [key: string]: number;
  }>({});
  const [isUploading, setIsUploading] = useState(false);
  const [preview, setPreview] = useState<string | null>(null);

  const uploadMutation = useUploadPhoto();

  // Generate a preview when files are selected
  const handleFileChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const selectedFiles = e.target.files;
      setFiles(selectedFiles);

      // Show preview of the first file
      if (selectedFiles && selectedFiles.length > 0) {
        const firstFile = selectedFiles[0];
        if (firstFile.type.startsWith("image/")) {
          const reader = new FileReader();
          reader.onloadend = () => {
            setPreview(reader.result as string);
          };
          reader.readAsDataURL(firstFile);
        } else {
          setPreview(null);
        }
      } else {
        setPreview(null);
      }
    },
    [],
  );

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!files || files.length === 0) return;

    try {
      setIsUploading(true);

      // Initialize progress for each file
      const initialProgress: { [key: string]: number } = {};
      for (let i = 0; i < files.length; i++) {
        initialProgress[files[i].name] = 0;
      }
      setUploadProgress(initialProgress);

      // Upload each file sequentially
      for (let i = 0; i < files.length; i++) {
        const file = files[i];
        const formData = new FormData();
        formData.append("file", file);

        // Update progress for this file
        setUploadProgress((prev) => ({ ...prev, [file.name]: 10 })); // Started

        // Upload file
        await uploadMutation.mutateAsync(formData);

        // Mark this file as complete
        setUploadProgress((prev) => ({ ...prev, [file.name]: 100 }));
      }

      // Reset form after successful upload
      setFiles(null);
      setPreview(null);
      setUploadProgress({});
      setIsUploading(false);
      onComplete();
    } catch (err) {
      console.error("Upload failed:", err);
      setIsUploading(false);
    }
  };

  // Calculate overall progress
  const calculateOverallProgress = (): number => {
    if (!files || files.length === 0) return 0;

    let total = 0;
    for (const filename in uploadProgress) {
      total += uploadProgress[filename];
    }
    return Math.round((total / (files.length * 100)) * 100);
  };

  return (
    <div className="upload-form">
      <h2>Upload Photos</h2>

      {uploadMutation.error && (
        <div className="error">
          {uploadMutation.error instanceof Error
            ? uploadMutation.error.message
            : "Upload failed"}
        </div>
      )}

      <form onSubmit={handleSubmit}>
        <div className={`file-input ${isUploading ? "disabled" : ""}`}>
          <label htmlFor="file-upload" className="custom-file-upload">
            <span className="upload-icon">📷</span>
            <span>Choose Photos</span>
          </label>
          <input
            id="file-upload"
            type="file"
            multiple
            accept="image/*"
            onChange={handleFileChange}
            disabled={isUploading}
          />

          {files && files.length > 0 && (
            <div className="selected-files">
              <p>{files.length} file(s) selected</p>
            </div>
          )}
        </div>

        {/* Image preview */}
        {preview && (
          <div className="upload-preview">
            <img src={preview} alt="Preview" />
            {files && files.length > 1 && (
              <p>+ {files.length - 1} more files</p>
            )}
          </div>
        )}

        {/* Progress bar (only show when uploading) */}
        {isUploading && (
          <div className="upload-progress">
            <div className="progress-bar">
              <div
                className="progress-fill"
                style={{ width: `${calculateOverallProgress()}%` }}
              ></div>
            </div>
            <p>{calculateOverallProgress()}% Complete</p>
          </div>
        )}

        <button
          type="submit"
          disabled={!files || files.length === 0 || isUploading}
          className="upload-button"
        >
          {isUploading
            ? "Uploading..."
            : `Upload ${files && files.length > 1 ? `${files.length} Photos` : "Photo"}`}
        </button>
      </form>
    </div>
  );
}
