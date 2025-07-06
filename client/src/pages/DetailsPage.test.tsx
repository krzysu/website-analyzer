import { render, screen } from "@testing-library/react";
import "@testing-library/jest-dom/vitest";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { useCrawlResultDetail } from "@/hooks/useCrawlResultDetail";
import { DetailsPage } from "@/pages/DetailsPage";
import {
  createLoadingMockUseQueryResult,
  createMockUseQueryResult,
  mockCrawlResult,
} from "@/test-utils/fixtures";

// Mock the useParams hook
vi.mock("react-router-dom", async () => {
  const actual = await vi.importActual("react-router-dom");
  return {
    ...actual,
    useParams: () => ({ id: "123" }), // Mock a default ID for testing
  };
});

// Mock the useCrawlResultDetail hook
vi.mock("@/hooks/useCrawlResultDetail");

// Mock the Pie component from react-chartjs-2
vi.mock("react-chartjs-2", () => ({
  Pie: vi.fn(() => null), // Mock Pie component to render nothing
}));

describe("DetailsPage", () => {
  // Get the mocked function and reset it before each test
  beforeEach(() => {
    vi.mocked(useCrawlResultDetail).mockReset();
  });

  it("renders loading state initially", () => {
    vi.mocked(useCrawlResultDetail).mockReturnValue({
      ...createLoadingMockUseQueryResult(undefined),
    });

    render(
      <MemoryRouter initialEntries={["/details/123"]}>
        <Routes>
          <Route path="/details/:id" element={<DetailsPage />} />
        </Routes>
      </MemoryRouter>,
    );

    expect(screen.getByText("Loading...")).toBeInTheDocument();
  });

  it("renders crawl result details when data is available", () => {
    vi.mocked(useCrawlResultDetail).mockReturnValue(
      createMockUseQueryResult(mockCrawlResult),
    );

    render(
      <MemoryRouter initialEntries={["/details/123"]}>
        <Routes>
          <Route path="/details/:id" element={<DetailsPage />} />
        </Routes>
      </MemoryRouter>,
    );

    expect(screen.getByText("Example Domain")).toBeInTheDocument();
    expect(screen.getByText("https://example.com")).toBeInTheDocument();
    expect(screen.getByText("HTML5")).toBeInTheDocument();
    expect(screen.getByText("No")).toBeInTheDocument(); // HasLoginForm
    expect(screen.getByText("Internal Links")).toBeInTheDocument();
    expect(screen.getByText("External Links")).toBeInTheDocument();
    expect(screen.getByText("H1")).toBeInTheDocument();
    expect(screen.getByText("1")).toBeInTheDocument();
    // More specific assertions for '2' to avoid ambiguity
    const h2Row = screen.getByRole("row", { name: /H2/i });
    expect(h2Row).toHaveTextContent("2");
    const externalLinksRow = screen.getByRole("row", {
      name: /External Links/i,
    });
    expect(externalLinksRow).toHaveTextContent("2");
    expect(screen.getByText(/No broken links found./i)).toBeInTheDocument();
  });
});
