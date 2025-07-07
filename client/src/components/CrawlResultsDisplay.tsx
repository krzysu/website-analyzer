import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { UrlTable } from "@/components/UrlTable";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { useCrawlResults } from "@/hooks/useCrawlResults";
import { useUrlActions } from "@/hooks/useUrlActions";

export function CrawlResultsDisplay() {
  const [selectedUrls, setSelectedUrls] = useState<number[]>([]);
  const [searchTerm, setSearchTerm] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 5;

  const navigate = useNavigate();

  const { crawlResults, totalResults } = useCrawlResults({
    polling: true,
    limit: itemsPerPage,
    offset: (currentPage - 1) * itemsPerPage,
    filterBy: searchTerm,
  });

  const totalPages = Math.ceil(totalResults / itemsPerPage);

  const paginate = (pageNumber: number) => setCurrentPage(pageNumber);

  const { handleBulkDelete: performBulkDelete, handleBulkRerun: performBulkRerun } =
    useUrlActions();

  const handleRowClick = (id: number) => {
    navigate(`/details/${id}`);
  };

  const handleCheckboxChange = (id: number) => {
    setSelectedUrls((prevSelected) =>
      prevSelected.includes(id)
        ? prevSelected.filter((urlId) => urlId !== id)
        : [...prevSelected, id],
    );
  };

  const handleBulkDelete = async () => {
    await performBulkDelete(selectedUrls);
    setSelectedUrls([]);
  };

  const handleBulkRerun = async () => {
    await performBulkRerun(selectedUrls);
    setSelectedUrls([]);
  };

  return (
    <div className="space-y-8">
      <Card className="w-full">
        <CardHeader className="flex flex-col md:flex-row justify-between items-start md:items-center gap-8">
          <CardTitle>Crawl Results</CardTitle>
          <div className="flex items-center space-x-4">
            <Button
              onClick={handleBulkDelete}
              disabled={selectedUrls.length === 0}
              size="sm"
              variant="destructive"
            >
              Delete Selected
            </Button>
            <Button
              onClick={handleBulkRerun}
              disabled={selectedUrls.length === 0}
              size="sm"
            >
              Re-run Selected
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <div className="mb-4">
            <Input
              type="text"
              placeholder="Search URLs, titles, or statuses..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="max-w-sm"
            />
          </div>
          <UrlTable
            results={crawlResults}
            onRowClick={handleRowClick}
            selectedUrls={selectedUrls}
            onCheckboxChange={handleCheckboxChange}
          />
          <Pagination className="mt-4">
            <PaginationContent>
              <PaginationItem>
                <PaginationPrevious
                  onClick={() => paginate(currentPage - 1)}
                  className={currentPage === 1 ? "pointer-events-none opacity-50" : ""}
                />
              </PaginationItem>
              {[...Array(totalPages)].map((_, i) => (
                <PaginationItem key={i + 1}>
                  <PaginationLink
                    onClick={() => paginate(i + 1)}
                    isActive={i + 1 === currentPage}
                  >
                    {i + 1}
                  </PaginationLink>
                </PaginationItem>
              ))}
              <PaginationItem>
                <PaginationNext
                  onClick={() => paginate(currentPage + 1)}
                  className={
                    currentPage === totalPages ? "pointer-events-none opacity-50" : ""
                  }
                />
              </PaginationItem>
            </PaginationContent>
          </Pagination>
        </CardContent>
      </Card>
    </div>
  );
}
