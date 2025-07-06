import { Route, Routes } from "react-router-dom";
import { DetailsPage } from "@/pages/DetailsPage";
import { HomePage } from "@/pages/HomePage";

export function App() {
  return (
    <div className="container mx-auto p-4">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold">Website Analyzer</h1>
      </div>

      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/details/:id" element={<DetailsPage />} />
      </Routes>
    </div>
  );
}
