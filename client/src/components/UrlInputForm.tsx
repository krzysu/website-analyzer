import type React from "react";
import { useState } from "react";

interface UrlInputFormProps {
  onSubmit: (url: string) => void;
}

const UrlInputForm: React.FC<UrlInputFormProps> = ({ onSubmit }) => {
  const [url, setUrl] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (url.trim()) {
      onSubmit(url);
      setUrl("");
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="text"
        value={url}
        onChange={(e) => setUrl(e.target.value)}
        placeholder="Enter URL to analyze"
      />
      <button type="submit">Analyze</button>
    </form>
  );
};

export default UrlInputForm;
