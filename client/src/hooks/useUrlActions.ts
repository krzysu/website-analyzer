import { useCallback } from "react";
import { useApi } from "./useApi";

export function useUrlActions(fetchCrawlResults: () => void) {
  const { callApi } = useApi();

  const handleUrlSubmit = useCallback(
    async (url: string) => {
      console.log("Submitting URL:", url);

      try {
        await callApi("/urls", {
          method: "POST",
          body: JSON.stringify({ url }),
        });

        console.log("URL submitted successfully:");
        fetchCrawlResults(); // Refresh the list after submitting a new URL
      } catch (error) {
        console.error("Error submitting URL:", error);
        // TODO: Display error message to user
      }
    },
    [callApi, fetchCrawlResults],
  );

  const handleBulkDelete = useCallback(
    async (selectedUrls: number[]) => {
      if (selectedUrls.length === 0) return;

      try {
        await callApi("/urls", {
          method: "DELETE",
          body: JSON.stringify({ ids: selectedUrls }),
        });

        console.log("URLs deleted successfully");
        fetchCrawlResults();
      } catch (error) {
        console.error("Error deleting URLs:", error);
      }
    },
    [callApi, fetchCrawlResults],
  );

  const handleBulkRerun = useCallback(
    async (selectedUrls: number[]) => {
      if (selectedUrls.length === 0) return;

      try {
        await callApi("/urls/rerun", {
          method: "POST",
          body: JSON.stringify({ ids: selectedUrls }),
        });

        console.log("Re-run initiated successfully");
        fetchCrawlResults();
      } catch (error) {
        console.error("Error re-running analysis:", error);
      }
    },
    [callApi, fetchCrawlResults],
  );

  return { handleUrlSubmit, handleBulkDelete, handleBulkRerun };
}
