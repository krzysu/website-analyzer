import type React from "react";
import type { CrawlResult } from "../types.ts";

interface UrlTableProps {
  results: CrawlResult[];
  onRowClick: (id: number) => void;
  selectedUrls: number[];
  onCheckboxChange: (id: number) => void;
}

const UrlTable: React.FC<UrlTableProps> = ({
  results,
  onRowClick,
  selectedUrls,
  onCheckboxChange,
}) => {
  return (
    <table>
      <thead>
        <tr>
          <th></th> {/* Checkbox column */}
          <th>Title</th>
          <th>HTML Version</th>
          <th>Internal Links</th>
          <th>External Links</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
        {results.map((result) => (
          <tr
            key={result.ID}
            onClick={() => onRowClick(result.ID)}
            onKeyDown={(e) => {
              if (e.key === "Enter" || e.key === " ") {
                onRowClick(result.ID);
              }
            }}
            tabIndex={0}
            style={{ cursor: "pointer" }}
          >
            <td>
              <input
                type="checkbox"
                checked={selectedUrls.includes(result.ID)}
                onChange={() => onCheckboxChange(result.ID)}
              />
            </td>
            <td style={{ cursor: "pointer" }}>
              {result.PageTitle || "Fetching title..."}
            </td>
            <td>{result.HTMLVersion}</td>
            <td>{result.InternalLinksCount}</td>
            <td>{result.ExternalLinksCount}</td>
            <td>{result.Status}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

export default UrlTable;
