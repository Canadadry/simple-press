import React, { useState } from "react";
import "./Uploader.css";

interface SingleFileUploaderProps {
  handleUpload: (f: Blob) => Promise<void>;
}

function SingleFileUploader({ handleUpload }: SingleFileUploaderProps) {
  const [file, setFile] = useState<File | null>(null);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      setFile(e.target.files[0]);
    }
  };

  return (
    <>
      <div className="input-group">
        <input id="file" type="file" onChange={handleFileChange} />
      </div>
      {file && (
        <section>
          File details:
          <ul>
            <li>Name: {file.name}</li>
            <li>Type: {file.type}</li>
            <li>Size: {file.size} bytes</li>
          </ul>
        </section>
      )}

      {file && (
        <button
          className="submit"
          onClick={() => {
            handleUpload(file);
            setFile(null);
          }}
        >
          Upload a file
        </button>
      )}
    </>
  );
}

export default SingleFileUploader;
