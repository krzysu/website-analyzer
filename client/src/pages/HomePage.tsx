import { useState } from "react";
import { CrawlResultsDisplay } from "@/components/CrawlResultsDisplay";
import { NoResultsDisplay } from "@/components/NoResultsDisplay";
import { useCrawlResults } from "@/hooks/useCrawlResults";
import { useUrlActions } from "@/hooks/useUrlActions";

export function HomePage() {
  const [isBulkDialogOpen, setBulkDialogOpen] = useState(false);

  const { crawlResults } = useCrawlResults();
  const { handleUrlSubmit, handleBulkUrlSubmit } = useUrlActions();

  return (
    <>
      {crawlResults.length === 0 ? (
        <NoResultsDisplay
          handleUrlSubmit={handleUrlSubmit}
          handleBulkUrlSubmit={handleBulkUrlSubmit}
          isBulkDialogOpen={isBulkDialogOpen}
          setBulkDialogOpen={setBulkDialogOpen}
        />
      ) : (
        <CrawlResultsDisplay />
      )}
    </>
  );
}
