import { useState } from "react";
import { Route, Routes, useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import UrlDetail from "./components/UrlDetail";
import UrlInputForm from "./components/UrlInputForm";
import UrlTable from "./components/UrlTable";
import { useCrawlResults } from "./hooks/useCrawlResults";
import { useUrlActions } from "./hooks/useUrlActions";

function App() {
  const [selectedUrls, setSelectedUrls] = useState<number[]>([]);
  const navigate = useNavigate();

  const { crawlResults, fetchCrawlResults } = useCrawlResults();
  const {
    handleUrlSubmit,
    handleBulkDelete: performBulkDelete,
    handleBulkRerun: performBulkRerun,
  } = useUrlActions(fetchCrawlResults);

  const handleRowClick = (id: number) => {
    navigate(`/details/${id}`);
  };

  const handleCheckboxChange = (id: number) => {
    setSelectedUrls((prevSelected) =>
      prevSelected.includes(id)
        ? prevSelected.filter((urlId) => urlId !== id)
        : [...prevSelected, id],
    );
  };

  const handleBulkDelete = async () => {
    await performBulkDelete(selectedUrls);
    setSelectedUrls([]);
  };

  const handleBulkRerun = async () => {
    await performBulkRerun(selectedUrls);
    setSelectedUrls([]);
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
                <Button
                  type="button"
                  onClick={handleBulkDelete}
                  disabled={selectedUrls.length === 0}
                >
                  Delete Selected
                </Button>
                <Button
                  type="button"
                  onClick={handleBulkRerun}
                  disabled={selectedUrls.length === 0}
                >
                  Re-run Selected
                </Button>
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
