import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { UrlInputForm } from "@/components/UrlInputForm";
import { UrlTable } from "@/components/UrlTable";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useCrawlResults } from "@/hooks/useCrawlResults";
import { useUrlActions } from "@/hooks/useUrlActions";

export function HomePage() {
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
    <>
      {crawlResults.length === 0 ? (
        <Card className="w-full">
          <CardHeader>
            <CardTitle>No Websites Analyzed Yet</CardTitle>
          </CardHeader>
          <CardContent className="text-center">
            <p className="mb-4 text-6xl">🚀</p>
            <p className="mb-4">Start by analyzing your first website.</p>
            <div className="flex justify-center">
              <UrlInputForm onSubmit={handleUrlSubmit} />
            </div>
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-4">
          <div className="space-y-4">
            <Card className="w-full">
              <CardHeader className="flex flex-row justify-between items-center">
                <CardTitle>Crawl Results</CardTitle>
              </CardHeader>
              <div className="flex items-center space-x-2 px-6 pb-4">
                <Button
                  onClick={handleBulkDelete}
                  disabled={selectedUrls.length === 0}
                  size="sm"
                >
                  Delete Selected
                </Button>
                <Button
                  onClick={handleBulkRerun}
                  disabled={selectedUrls.length === 0}
                  size="sm"
                >
                  Re-run Selected
                </Button>
              </div>
              <CardContent>
                <UrlTable
                  results={crawlResults}
                  onRowClick={handleRowClick}
                  selectedUrls={selectedUrls}
                  onCheckboxChange={handleCheckboxChange}
                />
              </CardContent>
            </Card>
          </div>
        </div>
      )}
    </>
  );
}
