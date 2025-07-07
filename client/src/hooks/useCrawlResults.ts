import { useQuery } from "@tanstack/react-query";
import type { CrawlResult } from "@/types.ts";
import { useApi } from "./useApi";

interface UseCrawlResultsOptions {
  polling?: boolean;
  limit?: number;
  offset?: number;
  sortBy?: string;
  filterBy?: string;
}

interface CrawlResultsResponse {
  results: CrawlResult[];
  total: number;
}

export function useCrawlResults({
  polling = false,
  limit = 10,
  offset = 0,
  sortBy = "created_at",
  filterBy = "",
}: UseCrawlResultsOptions = {}) {
  const { callApi } = useApi();

  const { data, ...queryInfo } = useQuery<CrawlResultsResponse>({
    queryKey: ["crawlResults", limit, offset, sortBy, filterBy],
    queryFn: () =>
      callApi<CrawlResultsResponse>(
        `/urls?limit=${limit}&offset=${offset}&sortBy=${sortBy}&filterBy=${filterBy}`,
      ),
    refetchInterval: polling ? 5000 : false,
  });

  return {
    crawlResults: data?.results || [],
    totalResults: data?.total || 0,
    ...queryInfo,
  };
}
