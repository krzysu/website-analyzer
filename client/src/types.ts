export interface CrawlResult {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  URL: string;
  Status: string;
  PageTitle: string;
  HTMLVersion: string;
  Headings: { [key: string]: number };
  InternalLinksCount: number;
  ExternalLinksCount: number;
  InaccessibleLinksCount: number;
  BrokenLinks: Array<{ url: string; statusCode: number }>;
  HasLoginForm: boolean;
  ErrorMessage: string;
}
