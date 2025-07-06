import { useQuery } from "@tanstack/react-query";
import type { CrawlResult } from "@/types.ts";
import { useApi } from "./useApi";

export function useCrawlResultDetail(id: string) {
  const { callApi } = useApi();

  return useQuery<CrawlResult | undefined>({
    queryKey: ["crawlResult", id],
    queryFn: () => callApi<CrawlResult>(`/urls/${id}`),
    enabled: !!id, // Only run the query if id is available
  });
}
