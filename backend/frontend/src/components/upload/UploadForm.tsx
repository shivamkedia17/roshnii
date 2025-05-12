import { useState, useCallback, useEffect, useRef } from "react";
import { useUploadImage } from "@/hooks/useImages";
import "@/css/Upload.css";

export type UploadFormProps = {
  onComplete: () => void;
};

export function UploadForm({ onComplete }: UploadFormProps) {
  const [files, setFiles] = useState<FileList | null>(null);
  const [uploadProgress, setUploadProgress] = useState<{
    [key: string]: number;
  }>({});
  const [isUploading, setIsUploading] = useState(false);
  const [preview, setPreview] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  // Get the upload mutation function from custom hook
  const uploadImage = useUploadImage();

  // Generate a preview when files are selected
  const handleFileChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const selectedFiles = e.target.files;
      setFiles(selectedFiles);
      setError(null);

      // Validate file types
      if (selectedFiles) {
        for (let i = 0; i < selectedFiles.length; i++) {
          if (!selectedFiles[i].type.startsWith("image/")) {
            setError("Please select image files only");
            return;
          }
        }
      }

      // Reset progress when new files are selected
      setUploadProgress({});

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

  // Handle form submission
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!files || files.length === 0) return;

    setIsUploading(true);
    setError(null);

    // Initialize progress for each file
    const initialProgress: { [key: string]: number } = {};
    for (let i = 0; i < files.length; i++) {
      initialProgress[files[i].name] = 0;
    }
    setUploadProgress(initialProgress);

    try {
      // Upload each file sequentially
      for (let i = 0; i < files.length; i++) {
        const file = files[i];

        // Update progress for this file - starting
        setUploadProgress((prev) => ({ ...prev, [file.name]: 10 }));

        // Upload the file using our hook
        uploadImage(file);

        // Mark this file as complete
        setUploadProgress((prev) => ({ ...prev, [file.name]: 100 }));
      }

      // Reset form after successful upload
      setFiles(null);
      setPreview(null);
      if (fileInputRef.current) {
        fileInputRef.current.value = "";
      }

      // Call the onComplete callback
      onComplete();
    } catch (error) {
      console.error("Error uploading files:", error);
      setError(
        error instanceof Error ? error.message : "Failed to upload images",
      );
    } finally {
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

  // Clean up preview URL when component unmounts or preview changes
  useEffect(() => {
    return () => {
      if (preview && preview.startsWith("blob:")) {
        URL.revokeObjectURL(preview);
      }
    };
  }, [preview]);

  return (
    <div className="upload-form">
      <h2>Upload Photos</h2>

      {error && <div className="error-message">{error}</div>}

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
            ref={fileInputRef}
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
            <p>
              {calculateOverallProgress()}% Complete
              {files &&
                files.length > 1 &&
                ` (${Object.keys(uploadProgress).filter((key) => uploadProgress[key] === 100).length}/${files.length})`}
            </p>
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
