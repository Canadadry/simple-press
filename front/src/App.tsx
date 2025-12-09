import { AppLayout } from "./AppLayout";
import Dashboard from "./pages/Dashboard";
import Articles from "./pages/articles-list/Articles";
import ArticleEdit from "./pages/articles-edit/ArticleEdit";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import Layouts from "./pages/layouts-list/Layouts";
import Templates from "./pages/templates-list/Templates";

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
          <Route path="/layouts" element={<Layouts />} />
          <Route path="/templates" element={<Templates />} />
          <Route path="/articles/:slug/edit" element={<ArticleEdit />} />
        </Routes>
      </AppLayout>
    </BrowserRouter>
  );
}
