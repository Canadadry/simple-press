import { AppLayout } from "./AppLayout";
import Dashboard from "./pages/Dashboard";
import Articles from "./pages/articles-list/Articles";
import ArticleEdit from "./pages/articles-edit/ArticleEdit";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import Layouts from "./pages/layouts-list/Layouts";
import Templates from "./pages/templates-list/Templates";
import ArticlePreview from "./pages/articles-preview/ArticlePreview";
import TemplateEdit from "./pages/templates-edit/TemplateEdit";

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
          <Route path="/templates/:slug/edit" element={<TemplateEdit />} />
          <Route path="/articles/:slug/edit" element={<ArticleEdit />} />
          <Route path="/articles/:slug/preview" element={<ArticlePreview />} />
        </Routes>
      </AppLayout>
    </BrowserRouter>
  );
}
