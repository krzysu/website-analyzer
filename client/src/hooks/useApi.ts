import { useCallback } from "react";

export function useApi() {
  const apiKey = import.meta.env.VITE_API_KEY;
  const apiBaseUrl = import.meta.env.VITE_API_BASE_URL;

  const callApi = useCallback(
    async <T>(endpoint: string, options?: RequestInit): Promise<T> => {
      try {
        const response = await fetch(`${apiBaseUrl}${endpoint}`, {
          ...options,
          headers: {
            "X-API-Key": apiKey,
            "Content-Type": "application/json",
            ...(options?.headers || {}),
          },
        });

        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.error || `API error: ${response.statusText}`);
        }

        return response.json();
      } catch (error) {
        console.error("API call error:", error);
        throw error;
      }
    },
    [],
  );

  return { callApi };
}
