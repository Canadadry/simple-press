import { AppLayout } from "./AppLayout";
import Dashboard from "./pages/Dashboard";
import Articles from "./pages/articles/Articles";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";

interface AppProps {
  toggleTheme: () => void;
}

export default function App({ toggleTheme }: AppProps) {
  return (
    <BrowserRouter>
      <AppLayout toggleTheme={toggleTheme}>
        <Routes>
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/articles" element={<Articles />} />
        </Routes>
      </AppLayout>
    </BrowserRouter>
  );
}
