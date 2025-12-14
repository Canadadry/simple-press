import { useEffect, useState } from "react";
import {
  Text,
  Flex,
  Spinner,
  Card,
  Box,
  Button,
  TextField,
} from "@radix-ui/themes";
// import * as Accordion from "@radix-ui/react-accordion";
import {
  getArticleList,
  postArticleAdd,
  type Article,
} from "../../api/article";
import Line from "./components/Line";
import { useNavigate } from "react-router-dom";

type SavingStatus = "untouched" | "touched" | "saving";

interface CreateProps {
  tabIndex: number;
}

function Create({ tabIndex }: CreateProps) {
  const [saving, setSaving] = useState<SavingStatus>("untouched");
  const [title, setTitle] = useState<string>("");
  const [author, setAuthor] = useState<string>("");
  const navigate = useNavigate();
  return (
    <Flex gap="3" mb="5">
      <Box flexGrow="3">
        <TextField.Root
          tabIndex={tabIndex}
          size="2"
          placeholder="Title"
          value={title}
          disabled={saving === "saving"}
          onChange={(e) => {
            setSaving("touched");
            setTitle(e.target.value);
          }}
        >
          <TextField.Slot>titre</TextField.Slot>
        </TextField.Root>
      </Box>
      <Box flexGrow="1">
        <TextField.Root
          tabIndex={tabIndex}
          size="2"
          placeholder="Author"
          value={author}
          disabled={saving === "saving"}
          onChange={(e) => {
            setSaving("touched");
            setAuthor(e.target.value);
          }}
        >
          <TextField.Slot>auteur</TextField.Slot>
        </TextField.Root>
      </Box>
      <Button
        tabIndex={tabIndex}
        size="2"
        disabled={saving != "touched"}
        onClick={async () => {
          setSaving("saving");
          await postArticleAdd({ title: title, author: author });
          navigate(`/articles/${title}/edit`, { replace: true });
          setSaving("untouched");
        }}
      >
        {saving == "saving" ? <Spinner /> : "Create"}
      </Button>
    </Flex>
  );
}

export default function Articles() {
  const [articles, setArticles] = useState<Article[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      try {
        const res = await getArticleList();
        setArticles(res.items);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  if (loading) {
    return (
      <Flex direction="column" gap="4">
        <Text size="7" weight="bold">
          Liste des articles
        </Text>
        <Create tabIndex={0}></Create>
        <Flex align="center" justify="center" height="100vh">
          <Spinner />
        </Flex>
      </Flex>
    );
  }

  return (
    <Flex direction="column" gap="4">
      <Text size="7" weight="bold">
        Liste des articles
      </Text>
      <Create tabIndex={0}></Create>
      <Card>
        <Flex direction="column">
          {articles.map((val, idx) => {
            return (
              <Line
                key={idx}
                tabIndex={idx}
                article={val}
                portalContainer={null}
              ></Line>
            );
          })}
        </Flex>
      </Card>
    </Flex>
  );
}
