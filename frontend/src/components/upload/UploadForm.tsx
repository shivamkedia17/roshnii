import { useState } from "react";
import { photosAPI } from "../../services/api";

interface UploadFormProps {
  onComplete: () => void;
}

export function UploadForm({ onComplete }: UploadFormProps) {
  const [file, setFile] = useState<File | null>(null);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!file) return;

    const formData = new FormData();
    formData.append("file", file);

    try {
      setUploading(true);
      await photosAPI.uploadPhoto(formData);
      setFile(null);
      onComplete();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Upload failed");
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="upload-form">
      <h2>Upload New Photo</h2>
      {error && <div className="error">{error}</div>}

      <form onSubmit={handleSubmit}>
        <div className="file-input">
          <input
            type="file"
            accept="image/*"
            onChange={(e) => setFile(e.target.files?.[0] || null)}
          />
        </div>

        <button
          type="submit"
          disabled={!file || uploading}
          className="upload-button"
        >
          {uploading ? "Uploading..." : "Upload Photo"}
        </button>
      </form>
    </div>
  );
}
