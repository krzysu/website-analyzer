import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { getStatusEmoji } from "@/lib/emojis";
import type { CrawlResult } from "@/types.ts";

interface UrlTableProps {
  results: CrawlResult[];
  onRowClick: (id: number) => void;
  selectedUrls: number[];
  onCheckboxChange: (id: number) => void;
}

export function UrlTable({
  results,
  onRowClick,
  selectedUrls,
  onCheckboxChange,
}: UrlTableProps) {
  return (
    <div className="relative w-full overflow-auto">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead title="Select for bulk actions">
              <span className="sr-only">Bulk Actions</span>
            </TableHead>
            <TableHead>URL</TableHead>
            <TableHead>Title</TableHead>
            <TableHead>Status</TableHead>
            <TableHead className="text-right">Details</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {results.map((result) => (
            <TableRow key={result.ID}>
              <TableCell onClick={(e) => e.stopPropagation()}>
                <Checkbox
                  checked={selectedUrls.includes(result.ID)}
                  onCheckedChange={() => onCheckboxChange(result.ID)}
                />
              </TableCell>
              <TableCell className="max-w-xs overflow-hidden text-ellipsis whitespace-nowrap">
                {result.URL}
              </TableCell>
              <TableCell className="max-w-xs overflow-hidden text-ellipsis whitespace-nowrap">
                {result.PageTitle || "Fetching title..."}
              </TableCell>
              <TableCell>
                <span>{getStatusEmoji(result.Status)}</span> <span>{result.Status}</span>
              </TableCell>
              <TableCell className="text-right">
                <Button variant="outline" size="sm" onClick={() => onRowClick(result.ID)}>
                  View
                </Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}
