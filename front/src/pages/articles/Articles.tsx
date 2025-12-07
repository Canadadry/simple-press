import { useEffect, useState } from "react";
import { Text, Flex, Spinner } from "@radix-ui/themes";
// import * as Accordion from "@radix-ui/react-accordion";
import { getArticleList } from "../../api/article";
import type { Article } from "../../api/article";
import Line from "./components/Line";

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
      <Flex align="center" justify="center" height="100vh">
        <Spinner />
      </Flex>
    );
  }

  return (
    <Flex direction="column" gap="4">
      <Text size="7" weight="bold">
        Liste des articles
      </Text>
      {articles.map((val, idx) => {
        return (
          <Line tabIndex={idx} article={val} portalContainer={null}></Line>
        );
      })}
    </Flex>
  );
}
