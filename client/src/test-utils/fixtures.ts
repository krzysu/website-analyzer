import type { QueryObserverLoadingResult, UseQueryResult } from "@tanstack/react-query";
import { vi } from "vitest";
import type { CrawlResult } from "@/types";

export const mockCrawlResult: CrawlResult = {
  ID: 1,
  URL: "https://example.com",
  PageTitle: "Example Domain",
  HTMLVersion: "HTML5",
  HasLoginForm: false,
  InternalLinksCount: 5,
  ExternalLinksCount: 2,
  InaccessibleLinksCount: 0,
  BrokenLinks: [],
  Headings: {
    H1: 1,
    H2: 2,
    H3: 0,
    H4: 0,
    H5: 0,
    H6: 0,
  },
  CreatedAt: "2023-01-01T12:00:00Z",
  UpdatedAt: "2023-01-01T12:00:00Z",
  ErrorMessage: "",
  Status: "completed",
};

export const mockCrawlResults: CrawlResult[] = [mockCrawlResult];

export function createMockUseQueryResult<T>(data: T): UseQueryResult<T, Error> {
  return {
    data,
    dataUpdatedAt: Date.now(),
    error: null,
    errorUpdatedAt: 0,
    failureCount: 0,
    failureReason: null,
    errorUpdateCount: 0,
    fetchStatus: "idle",

    // Status flags
    isError: false,
    isFetched: true,
    isFetchedAfterMount: true,
    isFetching: false,
    isInitialLoading: false,
    isLoading: false,
    isLoadingError: false,
    isPaused: false,
    isPending: false,
    isPlaceholderData: false,
    isRefetchError: false,
    isRefetching: false,
    isStale: false,
    isSuccess: true,
    status: "success",

    // Methods
    refetch: vi.fn(),

    // Deprecated but sometimes required
    promise: Promise.resolve(data),
  };
}

export function createLoadingMockUseQueryResult<T extends undefined>(
  data: T,
): QueryObserverLoadingResult<T, Error> {
  return {
    data,
    dataUpdatedAt: 0,
    error: null,
    errorUpdatedAt: 0,
    failureCount: 0,
    failureReason: null,
    errorUpdateCount: 0,
    fetchStatus: "fetching",

    // Status flags
    isError: false,
    isFetched: false,
    isFetchedAfterMount: false,
    isFetching: true,
    isInitialLoading: true,
    isLoading: true,
    isLoadingError: false,
    isPaused: false,
    isPending: true,
    isPlaceholderData: false,
    isRefetchError: false,
    isRefetching: false,
    isStale: false,
    isSuccess: false,
    status: "pending",

    // Methods
    refetch: vi.fn(),

    // Promise
    promise: Promise.resolve(data),
  };
}
