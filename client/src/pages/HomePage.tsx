import { useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import { UrlInputForms } from "@/components/UrlInputForms";
import { UrlTable } from "@/components/UrlTable";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { useCrawlResults } from "@/hooks/useCrawlResults";
import { useUrlActions } from "@/hooks/useUrlActions";

export function HomePage() {
  const [selectedUrls, setSelectedUrls] = useState<number[]>([]);
  const [isBulkDialogOpen, setBulkDialogOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState("");
  const navigate = useNavigate();

  const { crawlResults } = useCrawlResults({ polling: true });

  const filteredCrawlResults = useMemo(() => {
    return crawlResults.filter(
      (result) =>
        result.URL.toLowerCase().includes(searchTerm.toLowerCase()) ||
        result.PageTitle.toLowerCase().includes(searchTerm.toLowerCase()) ||
        result.Status.toLowerCase().includes(searchTerm.toLowerCase()),
    );
  }, [crawlResults, searchTerm]);

  const {
    handleUrlSubmit,
    handleBulkDelete: performBulkDelete,
    handleBulkRerun: performBulkRerun,
    handleBulkUrlSubmit,
  } = useUrlActions();

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
        <Card className="w-full max-w-lg mx-auto">
          <CardHeader>
            <CardTitle className="text-center">No Websites Analyzed Yet</CardTitle>
          </CardHeader>
          <CardContent className="text-center">
            <p className="mb-6 text-6xl">🚀</p>
            <p className="mb-6">Start by analyzing your first website.</p>
            <UrlInputForms
              handleUrlSubmit={handleUrlSubmit}
              handleBulkUrlSubmit={handleBulkUrlSubmit}
              isBulkDialogOpen={isBulkDialogOpen}
              setBulkDialogOpen={setBulkDialogOpen}
              autoFocus={true}
            />
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-8">
          <Card className="w-full">
            <CardHeader className="flex flex-col md:flex-row justify-between items-start md:items-center gap-8">
              <CardTitle>Crawl Results</CardTitle>
              <div className="flex items-center space-x-4">
                <Button
                  onClick={handleBulkDelete}
                  disabled={selectedUrls.length === 0}
                  size="sm"
                  variant="destructive"
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
            </CardHeader>
            <CardContent>
              <div className="mb-4">
                <Input
                  type="text"
                  placeholder="Search URLs, titles, or statuses..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="max-w-sm"
                />
              </div>
              <UrlTable
                results={filteredCrawlResults}
                onRowClick={handleRowClick}
                selectedUrls={selectedUrls}
                onCheckboxChange={handleCheckboxChange}
              />
            </CardContent>
          </Card>
        </div>
      )}
    </>
  );
}
