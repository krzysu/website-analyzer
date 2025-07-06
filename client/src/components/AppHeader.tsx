import { useLocation, useNavigate } from "react-router-dom";
import { UrlInputForm } from "@/components/UrlInputForm";
import { Button } from "@/components/ui/button";
import { useCrawlResults } from "@/hooks/useCrawlResults";
import { useUrlActions } from "@/hooks/useUrlActions";

export function AppHeader() {
  const location = useLocation();
  const navigate = useNavigate();
  const { crawlResults } = useCrawlResults();
  const { handleUrlSubmit } = useUrlActions();

  return (
    <div className="flex flex-col md:flex-row justify-between items-center mb-12 gap-8">
      <h1 className="text-3xl font-bold">Website Analyzer</h1>
      {location.pathname === "/" && crawlResults.length > 0 && (
        <div className="flex-grow flex justify-center md:justify-end w-full md:w-auto">
          <UrlInputForm onSubmit={handleUrlSubmit} />
        </div>
      )}
      {location.pathname.startsWith("/details/") && (
        <Button onClick={() => navigate("/")} variant="outline" size="sm">
          Back to Home
        </Button>
      )}
    </div>
  );
}
