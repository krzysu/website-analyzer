import { Route, Routes } from "react-router-dom";
import { AppHeader } from "@/components/AppHeader";
import { DetailsPage } from "@/pages/DetailsPage";
import { HomePage } from "@/pages/HomePage";

export function App() {
  return (
    <div className="max-w-7xl mx-auto p-8">
      <AppHeader />

      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/details/:id" element={<DetailsPage />} />
      </Routes>
    </div>
  );
}
