import { UrlInputForms } from "@/components/UrlInputForms";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface NoResultsDisplayProps {
  handleUrlSubmit: (url: string) => void;
  handleBulkUrlSubmit: (urls: string[]) => void;
  isBulkDialogOpen: boolean;
  setBulkDialogOpen: (open: boolean) => void;
}

export function NoResultsDisplay({
  handleUrlSubmit,
  handleBulkUrlSubmit,
  isBulkDialogOpen,
  setBulkDialogOpen,
}: NoResultsDisplayProps) {
  return (
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
  );
}
