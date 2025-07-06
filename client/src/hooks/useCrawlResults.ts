import { useQuery } from "@tanstack/react-query";
import type { CrawlResult } from "@/types.ts";
import { useApi } from "./useApi";

export function useCrawlResults({ polling = false }: { polling?: boolean } = {}) {
  const { callApi } = useApi();

  const { data: crawlResults, ...queryInfo } = useQuery<CrawlResult[]>({
    queryKey: ["crawlResults"],
    queryFn: () => callApi<CrawlResult[]>("/urls"),
    refetchInterval: polling ? 5000 : false,
  });

  return { crawlResults: crawlResults || [], ...queryInfo };
}
