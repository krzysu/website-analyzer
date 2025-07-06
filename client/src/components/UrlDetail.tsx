import { ArcElement, Chart as ChartJS, Legend, Tooltip } from "chart.js";
import type React from "react";
import { useEffect, useState } from "react";
import { Pie } from "react-chartjs-2";
import { useParams } from "react-router-dom";
import type { CrawlResult } from "../types.ts";

ChartJS.register(ArcElement, Tooltip, Legend);

const UrlDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [crawlResult, setCrawlResult] = useState<CrawlResult | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchUrlDetail = async () => {
      const apiKey = import.meta.env.VITE_API_KEY;
      try {
        const response = await fetch(
          `${import.meta.env.VITE_API_BASE_URL}/urls/${id}`,
          {
            headers: {
              "X-API-Key": apiKey,
            },
          },
        );
        if (!response.ok) {
          throw new Error("Failed to fetch URL details");
        }
        const data: CrawlResult = await response.json();
        setCrawlResult(data);
      } catch (err) {
        setError((err as Error).message);
      } finally {
        setLoading(false);
      }
    };

    fetchUrlDetail();
  }, [id]);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  if (!crawlResult) {
    return <div>No data found.</div>;
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
    <div>
      <h2>Details for: {crawlResult.URL}</h2>
      <p>
        <strong>Status:</strong> {crawlResult.Status}
      </p>
      <p>
        <strong>Page Title:</strong> {crawlResult.PageTitle}
      </p>
      <p>
        <strong>HTML Version:</strong> {crawlResult.HTMLVersion}
      </p>
      <p>
        <strong>Has Login Form:</strong>{" "}
        {crawlResult.HasLoginForm ? "Yes" : "No"}
      </p>

      <h3>Heading Counts:</h3>
      <ul>
        {Object.entries(crawlResult.Headings).map(([heading, count]) => (
          <li key={heading}>
            {heading}: {count}
          </li>
        ))}
      </ul>

      <h3>Link Distribution:</h3>
      <div style={{ width: "300px", height: "300px" }}>
        <Pie data={chartData} />
      </div>

      <h3>Broken Links ({crawlResult.InaccessibleLinksCount}):</h3>
      {crawlResult.BrokenLinks.length > 0 ? (
        <ul>
          {crawlResult.BrokenLinks.map((link, _index) => (
            <li key={link.url}>
              {link.url} (Status: {link.statusCode})
            </li>
          ))}
        </ul>
      ) : (
        <p>No broken links found.</p>
      )}
    </div>
  );
};

export default UrlDetail;
