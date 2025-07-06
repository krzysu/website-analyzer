import { useState, useEffect } from "react";
import { Routes, Route, useNavigate } from "react-router-dom";
import UrlInputForm from "./components/UrlInputForm";
import UrlTable from "./components/UrlTable";
import UrlDetail from "./components/UrlDetail";
import type { CrawlResult } from "./types.ts";

function App() {
  const [crawlResults, setCrawlResults] = useState<CrawlResult[]>([]);
  const [selectedUrls, setSelectedUrls] = useState<number[]>([]);
  const navigate = useNavigate();

  const fetchCrawlResults = async () => {
    const apiKey = import.meta.env.VITE_API_KEY;
    try {
      const response = await fetch(
        `${import.meta.env.VITE_API_BASE_URL}/urls`,
        {
          headers: {
            "X-API-Key": apiKey,
          },
        }
      );
      if (!response.ok) {
        throw new Error("Failed to fetch crawl results");
      }
      const data: CrawlResult[] = await response.json();
      setCrawlResults(data);
    } catch (error) {
      console.error("Error fetching crawl results:", error);
    }
  };

  useEffect(() => {
    fetchCrawlResults();
    const intervalId = setInterval(fetchCrawlResults, 5000); // Poll every 5 seconds
    return () => clearInterval(intervalId);
  }, []);

  const handleUrlSubmit = async (url: string) => {
    console.log("Submitting URL:", url);
    const apiKey = import.meta.env.VITE_API_KEY;

    try {
      const response = await fetch(
        `${import.meta.env.VITE_API_BASE_URL}/urls`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            "X-API-Key": apiKey,
          },
          body: JSON.stringify({ url }),
        }
      );

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || "Failed to submit URL");
      }

      const data = await response.json();
      console.log("URL submitted successfully:", data);
      fetchCrawlResults(); // Refresh the list after submitting a new URL
    } catch (error) {
      console.error("Error submitting URL:", error);
      // TODO: Display error message to user
    }
  };

  const handleRowClick = (id: number) => {
    navigate(`/details/${id}`);
  };

  const handleCheckboxChange = (id: number) => {
    setSelectedUrls((prevSelected) =>
      prevSelected.includes(id)
        ? prevSelected.filter((urlId) => urlId !== id)
        : [...prevSelected, id]
    );
  };

  const handleBulkDelete = async () => {
    if (selectedUrls.length === 0) return;
    const apiKey = import.meta.env.VITE_API_KEY;

    try {
      const response = await fetch(
        `${import.meta.env.VITE_API_BASE_URL}/urls`,
        {
          method: "DELETE",
          headers: {
            "Content-Type": "application/json",
            "X-API-Key": apiKey,
          },
          body: JSON.stringify({ ids: selectedUrls }),
        }
      );

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || "Failed to delete URLs");
      }

      console.log("URLs deleted successfully");
      setSelectedUrls([]);
      fetchCrawlResults();
    } catch (error) {
      console.error("Error deleting URLs:", error);
    }
  };

  const handleBulkRerun = async () => {
    if (selectedUrls.length === 0) return;
    const apiKey = import.meta.env.VITE_API_KEY;

    try {
      const response = await fetch(
        `${import.meta.env.VITE_API_BASE_URL}/urls/rerun`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            "X-API-Key": apiKey,
          },
          body: JSON.stringify({ ids: selectedUrls }),
        }
      );

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || "Failed to re-run analysis");
      }

      console.log("Re-run initiated successfully");
      setSelectedUrls([]);
      fetchCrawlResults();
    } catch (error) {
      console.error("Error re-running analysis:", error);
    }
  };

  return (
    <div className="App">
      <h1 className="text-3xl font-bold">Website Analyzer</h1>
      <Routes>
        <Route
          path="/"
          element={
            <>
              <UrlInputForm onSubmit={handleUrlSubmit} />
              <div>
                <button
                  onClick={handleBulkDelete}
                  disabled={selectedUrls.length === 0}
                >
                  Delete Selected
                </button>
                <button
                  onClick={handleBulkRerun}
                  disabled={selectedUrls.length === 0}
                >
                  Re-run Selected
                </button>
              </div>
              <UrlTable
                results={crawlResults}
                onRowClick={handleRowClick}
                selectedUrls={selectedUrls}
                onCheckboxChange={handleCheckboxChange}
              />
            </>
          }
        />
        <Route path="/details/:id" element={<UrlDetail />} />
      </Routes>
    </div>
  );
}

export default App;
