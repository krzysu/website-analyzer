import { ArcElement, Chart as ChartJS, Legend, Tooltip } from "chart.js";
import { Pie } from "react-chartjs-2";
import { useParams } from "react-router-dom";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Table, TableBody, TableCell, TableRow } from "@/components/ui/table";
import { useCrawlResultDetail } from "@/hooks/useCrawlResultDetail";

ChartJS.register(ArcElement, Tooltip, Legend);

export function DetailsPage() {
  const { id } = useParams<{ id: string }>();
  const { data: crawlResult, isLoading, error } = useCrawlResultDetail(id || "");

  if (isLoading) {
    return <div className="text-center py-8">Loading...</div>;
  }

  if (error) {
    return <div className="text-center py-8 text-red-500">Error: {error.message}</div>;
  }

  if (!crawlResult) {
    return <div className="text-center py-8">No data found.</div>;
  }

  const chartData = {
    labels: ["Internal Links", "External Links"],
    datasets: [
      {
        data: [crawlResult.InternalLinksCount, crawlResult.ExternalLinksCount],
        backgroundColor: ["#36A2EB", "#FF6384"],
        hoverBackgroundColor: ["#36A2EB", "#FF6384"],
      },
    ],
  };

  return (
    <>
      <Card className="mb-6">
        <CardHeader className="flex flex-row justify-between items-center">
          <CardTitle>{crawlResult.PageTitle || crawlResult.URL}</CardTitle>
        </CardHeader>
        <CardContent className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div>
            <Label>URL:</Label>
            <p>
              <a
                href={crawlResult.URL}
                target="_blank"
                rel="noopener noreferrer"
                className="text-blue-500 hover:underline"
              >
                {crawlResult.URL}
              </a>
            </p>
          </div>
          <div>
            <Label>HTML Version:</Label>
            <p>{crawlResult.HTMLVersion || "N/A"}</p>
          </div>
          <div>
            <Label>Has Login Form:</Label>
            <p>{crawlResult.HasLoginForm ? "Yes" : "No"}</p>
          </div>
          <div>
            <Label>Created At:</Label>
            <p>{new Date(crawlResult.CreatedAt).toLocaleString()}</p>
          </div>
          <div>
            <Label>Updated At:</Label>
            <p>{new Date(crawlResult.UpdatedAt).toLocaleString()}</p>
          </div>
          {crawlResult.ErrorMessage && (
            <div className="lg:col-span-3">
              <Label>Error Message:</Label>
              <p className="text-red-500">{crawlResult.ErrorMessage}</p>
            </div>
          )}
        </CardContent>
      </Card>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Heading Counts</CardTitle>
          </CardHeader>
          <CardContent>
            {Object.keys(crawlResult.Headings).length > 0 ? (
              <Table>
                <TableBody>
                  {Object.entries(crawlResult.Headings).map(([heading, count]) => (
                    <TableRow key={heading}>
                      <TableCell className="font-medium">{heading}</TableCell>
                      <TableCell>{count}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            ) : (
              <p>No headings found.</p>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Link Distribution</CardTitle>
          </CardHeader>
          <CardContent className="flex flex-col items-center">
            {crawlResult.InternalLinksCount > 0 || crawlResult.ExternalLinksCount > 0 ? (
              <div className="mb-4" style={{ width: "300px", height: "300px" }}>
                <Pie data={chartData} />
              </div>
            ) : (
              <p>No links found.</p>
            )}
            <Table>
              <TableBody>
                <TableRow>
                  <TableCell className="font-medium">Internal Links</TableCell>
                  <TableCell>{crawlResult.InternalLinksCount}</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell className="font-medium">External Links</TableCell>
                  <TableCell>{crawlResult.ExternalLinksCount}</TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </div>

      <Card className="mt-6">
        <CardHeader>
          <CardTitle>Broken Links ({crawlResult.InaccessibleLinksCount})</CardTitle>
        </CardHeader>
        <CardContent>
          {crawlResult.BrokenLinks.length > 0 ? (
            <Table>
              <TableBody>
                {crawlResult.BrokenLinks.map((link, _index) => (
                  <TableRow key={link.url}>
                    <TableCell className="font-medium">{link.url}</TableCell>
                    <TableCell>{link.statusCode}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          ) : (
            <p>No broken links found.</p>
          )}
        </CardContent>
      </Card>
    </>
  );
}
