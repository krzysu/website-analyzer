import { useLocation } from "react-router-dom";
import { UrlInputForm } from "@/components/UrlInputForm";
import { useCrawlResults } from "@/hooks/useCrawlResults";
import { useUrlActions } from "@/hooks/useUrlActions";

export function AppHeader() {
  const location = useLocation();
  const { crawlResults } = useCrawlResults();
  const { handleUrlSubmit } = useUrlActions();

  return (
    <div className="flex justify-between items-center mb-8">
      <h1 className="text-3xl font-bold">Website Analyzer</h1>
      {location.pathname === "/" && crawlResults.length > 0 && (
        <div className="flex-grow flex justify-end">
          <UrlInputForm onSubmit={handleUrlSubmit} />
        </div>
      )}
    </div>
  );
}
