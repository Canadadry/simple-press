// src/pages/Articles.tsx
import { useEffect, useState } from "react";
import { Text, Flex, Card, Spinner } from "@radix-ui/themes";
import * as Accordion from "@radix-ui/react-accordion";
import { getArticleList } from "../api/article";
import type { Article } from "../api/article";

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

      <Card>
        <Accordion.Root type="multiple">
          {articles.map((article) => (
            <Accordion.Item key={article.slug} value={article.slug}>
              <Accordion.Header>
                <Accordion.Trigger>{article.title}</Accordion.Trigger>
              </Accordion.Header>

              <Accordion.Content>
                <Flex direction="column" gap="2">
                  <Text>
                    <b>Auteur :</b> {article.author}
                  </Text>
                  <Text>
                    <b>Slug :</b> {article.slug}
                  </Text>

                  {article.draft && <Text color="red">Brouillon</Text>}

                  {article.blocks && article.blocks.length > 0 && (
                    <Accordion.Root type="multiple">
                      <Text weight="bold" size="3">
                        Blocks
                      </Text>

                      {article.blocks.map((block, i) => (
                        <Accordion.Item
                          key={i}
                          value={`${article.slug}-block-${i}`}
                        >
                          <Accordion.Header>
                            <Accordion.Trigger>{block.name}</Accordion.Trigger>
                          </Accordion.Header>
                          <Accordion.Content>
                            <Text>Valeur : {block.value}</Text>
                          </Accordion.Content>
                        </Accordion.Item>
                      ))}
                    </Accordion.Root>
                  )}
                </Flex>
              </Accordion.Content>
            </Accordion.Item>
          ))}
        </Accordion.Root>
      </Card>
    </Flex>
  );
}
