import type { CrawlStatus } from "@/types";

export const getStatusEmoji = (status: CrawlStatus) => {
  switch (status) {
    case "completed":
      return "✅";
    case "error":
      return "❌";
    case "queued":
      return "⏳";
    case "running":
      return "⚙️";
    default:
      return "❓";
  }
};
