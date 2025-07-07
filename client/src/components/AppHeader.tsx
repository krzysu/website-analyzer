import { useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { BulkUrlInputForm } from "@/components/BulkUrlInputForm";
import { UrlInputForm } from "@/components/UrlInputForm";
import { UrlInputForms } from "@/components/UrlInputForms";
import { Button } from "@/components/ui/button";
import { useCrawlResults } from "@/hooks/useCrawlResults";
import { useUrlActions } from "@/hooks/useUrlActions";

export function AppHeader() {
  const [isBulkDialogOpen, setBulkDialogOpen] = useState(false);
  const location = useLocation();
  const navigate = useNavigate();
  const { crawlResults } = useCrawlResults();
  const { handleUrlSubmit, handleBulkUrlSubmit } = useUrlActions();

  return (
    <div className="flex flex-col md:flex-row justify-between items-center mb-12 gap-8">
      <h1 className="text-3xl font-bold">Website Analyzer</h1>

      {location.pathname === "/" && crawlResults.length > 0 && (
        <>
          <div className="hidden md:flex flex-grow justify-center md:justify-end w-full md:w-auto">
            <div className="flex items-center gap-2">
              <UrlInputForm onSubmit={handleUrlSubmit} className="w-full" />
              <BulkUrlInputForm
                onSubmit={handleBulkUrlSubmit}
                open={isBulkDialogOpen}
                onOpenChange={setBulkDialogOpen}
                displayMode="icon"
              />
            </div>
          </div>

          <div className="md:hidden w-full">
            <UrlInputForms
              handleUrlSubmit={handleUrlSubmit}
              handleBulkUrlSubmit={handleBulkUrlSubmit}
              isBulkDialogOpen={isBulkDialogOpen}
              setBulkDialogOpen={setBulkDialogOpen}
            />
          </div>
        </>
      )}

      {location.pathname.startsWith("/details/") && (
        <Button onClick={() => navigate("/")} variant="outline" size="sm">
          Back to Home
        </Button>
      )}
    </div>
  );
}
