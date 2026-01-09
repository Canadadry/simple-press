import { useEffect, useState } from "react";
import { Dict } from "../api/api";
import { Editor } from "./Editor";

interface JsonEditorProps {
  tabIndex: number;
  title: string;
  content: Dict;
  setContent: (content: Dict) => void;
  updateContent: () => Promise<void>;
}

export const JsonEditor: React.FC<JsonEditorProps> = ({
  tabIndex,
  title,
  content,
  setContent,
  updateContent,
}: JsonEditorProps) => {
  const [temp, setTemp] = useState<string>("{}");
  useEffect(() => {
    setTemp(JSON.stringify(content, null, 2));
  }, [content]);
  return (
    <Editor
      tabIndex={tabIndex}
      title={title}
      content={temp}
      setContent={(content: string) => {
        let p: Dict | null = null;
        try {
          p = JSON.parse(content);
        } catch {
          setTemp(content);
        } finally {
          if (p != null) {
            setTemp("{}");
            setContent(p);
          }
        }
      }}
      updateContent={updateContent}
    ></Editor>
  );
};
