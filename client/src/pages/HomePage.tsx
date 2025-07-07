import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { BulkUrlInputForm } from "@/components/BulkUrlInputForm";
import { UrlInputForm } from "@/components/UrlInputForm";
import { UrlTable } from "@/components/UrlTable";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useCrawlResults } from "@/hooks/useCrawlResults";
import { useUrlActions } from "@/hooks/useUrlActions";

export function HomePage() {
  const [selectedUrls, setSelectedUrls] = useState<number[]>([]);
  const [isBulkDialogOpen, setBulkDialogOpen] = useState(false);
  const navigate = useNavigate();

  const { crawlResults } = useCrawlResults({ polling: true });
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
            <div className="w-full max-w-md mx-auto space-y-4">
              <UrlInputForm onSubmit={handleUrlSubmit} autoFocus={true} />
              <div className="relative">
                <div className="absolute inset-0 flex items-center">
                  <span className="w-full border-t" />
                </div>
                <div className="relative flex justify-center text-xs uppercase">
                  <span className="bg-card px-2 text-muted-foreground">Or</span>
                </div>
              </div>
              <BulkUrlInputForm
                onSubmit={handleBulkUrlSubmit}
                open={isBulkDialogOpen}
                onOpenChange={setBulkDialogOpen}
                displayMode="full"
              />
            </div>
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
              <UrlTable
                results={crawlResults}
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
