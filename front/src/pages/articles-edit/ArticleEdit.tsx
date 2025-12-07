import { useEffect, useState } from "react";
import { Text, Flex, Spinner, Card } from "@radix-ui/themes";
// import * as Accordion from "@radix-ui/react-accordion";
import { getArticleEdit } from "../../api/article";
import type { Article } from "../../api/article";
import { useParams } from "react-router-dom";

export default function Articles() {
  const { slug } = useParams<{ slug: string }>();
  const [article, setArticle] = useState<Article | null>(null);

  useEffect(() => {
    async function load() {
      if (!slug) {
        setArticle(null);
        return;
      }
      const res = await getArticleEdit(slug);
      setArticle(res);
    }
    load();
  }, [slug]);

  if (!article) {
    return (
      <Flex align="center" justify="center" height="100vh">
        <Spinner />
      </Flex>
    );
  }

  return (
    <Flex direction="column" gap="4">
      <Text size="7" weight="bold">
        {article.title}
      </Text>
      <Card>
        <Flex direction="column"></Flex>
      </Card>
    </Flex>
  );
}
