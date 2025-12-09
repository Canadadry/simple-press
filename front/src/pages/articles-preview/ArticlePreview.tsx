import { Link } from "@radix-ui/themes";
import { Pencil2Icon } from "@radix-ui/react-icons";
import { useEffect, useState } from "react";
import { Text, Flex, Spinner, Card } from "@radix-ui/themes";
import { getArticleEdit } from "../../api/article";
import type { Article } from "../../api/article";
import { useNavigate, useParams } from "react-router-dom";

export default function ArticlePreview() {
  const navigate = useNavigate();
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
        <Link
          href="#"
          mx={"2"}
          onClick={(e) => {
            e.preventDefault();
            navigate(`/articles/${slug}/edit`, { replace: true });
          }}
        >
          <Pencil2Icon color={"#000"} width={20} height={20}></Pencil2Icon>
        </Link>
      </Text>
      <Card>
        <iframe
          style={{
            width: "100%",
            height: "400px",
          }}
          title={`preview of ${article.title}`}
          src={`http://localhost:8080/admin/articles/${slug}/preview`}
        ></iframe>
      </Card>
    </Flex>
  );
}
