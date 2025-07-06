import { useCallback, useEffect, useState } from "react";
import type { CrawlResult } from "../types.ts";
import { useApi } from "./useApi";

export function useCrawlResults() {
  const [crawlResults, setCrawlResults] = useState<CrawlResult[]>([]);
  const { callApi } = useApi();

  const fetchCrawlResults = useCallback(async () => {
    try {
      const data: CrawlResult[] = await callApi<CrawlResult[]>("/urls");
      setCrawlResults(data);
    } catch (error) {
      console.error("Error fetching crawl results:", error);
    }
  }, [callApi]);

  useEffect(() => {
    fetchCrawlResults();
    const intervalId = setInterval(fetchCrawlResults, 5000); // Poll every 5 seconds
    return () => clearInterval(intervalId);
  }, [fetchCrawlResults]);

  return { crawlResults, fetchCrawlResults };
}
