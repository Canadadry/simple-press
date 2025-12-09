import { Box, IconButton, Button, TextArea, Link } from "@radix-ui/themes";
import { TextField } from "@radix-ui/themes";
import {
  CrumpledPaperIcon,
  FontBoldIcon,
  FontItalicIcon,
  ImageIcon,
  MagicWandIcon,
  EyeOpenIcon,
  StrikethroughIcon,
  TextAlignCenterIcon,
  TextAlignLeftIcon,
  TextAlignRightIcon,
} from "@radix-ui/react-icons";
import { Label } from "@radix-ui/react-label";

import { useEffect, useState } from "react";
import { Text, Flex, Spinner, Card } from "@radix-ui/themes";
import {
  getArticleEdit,
  postArticleEditMetadata,
  postArticleEditContent,
} from "../../api/article";
import type { Article } from "../../api/article";
import { useNavigate, useParams } from "react-router-dom";

type SavingStatus = "untouched" | "touched" | "saving";
interface Saving {
  metadada: SavingStatus;
  content: SavingStatus;
}

export default function Articles() {
  const navigate = useNavigate();
  const { slug } = useParams<{ slug: string }>();
  const [article, setArticle] = useState<Article | null>(null);
  const [saving, setSaving] = useState<Saving>({
    metadada: "untouched",
    content: "untouched",
  });
  const tabIndex = 1;
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
            navigate(`/articles/${slug}/preview`, { replace: true });
          }}
        >
          <EyeOpenIcon color={"#000"} width={20} height={20}></EyeOpenIcon>
        </Link>
      </Text>
      <Card>
        <Flex direction="column">
          <Flex gap="3" mb="5">
            <Box flexGrow="2">
              <TextField.Root
                tabIndex={tabIndex}
                size="2"
                placeholder="Title"
                value={article.title}
                disabled={saving.metadada === "saving"}
                onChange={(e) => {
                  setSaving({ ...saving, metadada: "touched" });
                  setArticle({ ...article, title: e.target.value });
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
                value={article.author}
                disabled={saving.metadada === "saving"}
                onChange={(e) => {
                  setSaving({ ...saving, metadada: "touched" });
                  setArticle({ ...article, author: e.target.value });
                }}
              >
                <TextField.Slot>auteur</TextField.Slot>
              </TextField.Root>
            </Box>
            <Button
              tabIndex={tabIndex}
              size="2"
              disabled={saving.metadada != "touched"}
              onClick={async () => {
                setSaving({ ...saving, metadada: "saving" });
                await postArticleEditMetadata(article.slug, article);
                setSaving({ ...saving, metadada: "untouched" });
              }}
            >
              {saving.metadada == "saving" ? <Spinner /> : "Save"}
            </Button>
          </Flex>
          <Box mb="4">
            <Label htmlFor="skirt-description">
              <Text size="2" weight="bold" mb="2" asChild>
                <Box display="inline-block">Content</Box>
              </Text>
            </Label>
            <Box position="relative">
              <TextArea
                tabIndex={tabIndex}
                spellCheck={false}
                id="skirt-description"
                variant="soft"
                rows={10}
                defaultValue={article.content || ""}
                style={{ paddingTop: 48 }}
                value={article.content}
                disabled={saving.content === "saving"}
                onChange={(e) => {
                  setSaving({ ...saving, content: "touched" });
                  setArticle({ ...article, content: e.target.value });
                }}
              />
              <Box position="absolute" m="2" top="0" left="0" right="0">
                <Flex gap="4">
                  <Flex gap="1">
                    <IconButton tabIndex={tabIndex} variant="soft" highContrast>
                      <FontItalicIcon />
                    </IconButton>

                    <IconButton tabIndex={tabIndex} variant="soft" highContrast>
                      <FontBoldIcon />
                    </IconButton>

                    <IconButton tabIndex={tabIndex} variant="soft" highContrast>
                      <StrikethroughIcon />
                    </IconButton>
                  </Flex>

                  <Flex gap="1">
                    <IconButton tabIndex={tabIndex} variant="soft" highContrast>
                      <TextAlignLeftIcon />
                    </IconButton>

                    <IconButton tabIndex={tabIndex} variant="soft" highContrast>
                      <TextAlignCenterIcon />
                    </IconButton>

                    <IconButton tabIndex={tabIndex} variant="soft" highContrast>
                      <TextAlignRightIcon />
                    </IconButton>
                  </Flex>

                  <Flex gap="1">
                    <IconButton tabIndex={tabIndex} variant="soft" highContrast>
                      <MagicWandIcon />
                    </IconButton>

                    <IconButton tabIndex={tabIndex} variant="soft" highContrast>
                      <ImageIcon />
                    </IconButton>

                    <IconButton tabIndex={tabIndex} variant="soft" highContrast>
                      <CrumpledPaperIcon />
                    </IconButton>
                  </Flex>
                  <Flex gap="1">
                    <Button
                      tabIndex={tabIndex}
                      size="2"
                      disabled={saving.content != "touched"}
                      onClick={async () => {
                        setSaving({ ...saving, content: "saving" });
                        await postArticleEditContent(
                          article.slug,
                          article.content || "",
                        );
                        setSaving({ ...saving, content: "untouched" });
                      }}
                    >
                      {saving.content == "saving" ? <Spinner /> : "Save"}
                    </Button>
                  </Flex>
                </Flex>
              </Box>
            </Box>
          </Box>
        </Flex>
      </Card>
    </Flex>
  );
}
