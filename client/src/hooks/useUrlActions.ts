import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useApi } from "./useApi";

export function useUrlActions() {
  const { callApi } = useApi();
  const queryClient = useQueryClient();

  const { mutate: handleUrlSubmit } = useMutation({
    mutationFn: async (url: string) => {
      await callApi("/urls", {
        method: "POST",
        body: JSON.stringify({ url }),
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["crawlResults"] });
    },
  });

  const { mutate: handleBulkDelete } = useMutation({
    mutationFn: async (selectedUrls: number[]) => {
      if (selectedUrls.length === 0) return;
      await callApi("/urls", {
        method: "DELETE",
        body: JSON.stringify({ ids: selectedUrls }),
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["crawlResults"] });
    },
  });

  const { mutate: handleBulkRerun } = useMutation({
    mutationFn: async (selectedUrls: number[]) => {
      if (selectedUrls.length === 0) return;
      await callApi("/urls/rerun", {
        method: "POST",
        body: JSON.stringify({ ids: selectedUrls }),
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["crawlResults"] });
    },
  });

  return { handleUrlSubmit, handleBulkDelete, handleBulkRerun };
}
