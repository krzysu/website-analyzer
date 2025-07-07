import { BulkUrlInputForm } from "@/components/BulkUrlInputForm";
import { UrlInputForm } from "@/components/UrlInputForm";

interface UrlInputFormsProps {
  handleUrlSubmit: (url: string) => void;
  handleBulkUrlSubmit: (urls: string[]) => void;
  isBulkDialogOpen: boolean;
  setBulkDialogOpen: (open: boolean) => void;
  autoFocus?: boolean;
}

export function UrlInputForms({
  handleUrlSubmit,
  handleBulkUrlSubmit,
  isBulkDialogOpen,
  setBulkDialogOpen,
  autoFocus,
}: UrlInputFormsProps) {
  return (
    <div className="w-full space-y-4">
      <UrlInputForm onSubmit={handleUrlSubmit} autoFocus={autoFocus} />
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
  );
}
