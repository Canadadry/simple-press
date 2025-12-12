import {
  Box,
  IconButton,
  Button,
  TextArea,
  Link,
  Select,
} from "@radix-ui/themes";
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
  postArticleEditBlockAdd,
} from "../../api/article";
import type { Article } from "../../api/article";
import { useNavigate, useParams } from "react-router-dom";
import { DynamicForm } from "../../pkg/data/render";
import { makeRadixUI } from "../../pkg/data/radix-form";

type SavingStatus = "untouched" | "touched" | "saving";

interface MetadataProps {
  tabIndex: number;
  article: Article;
  setArticle: (a: Article) => void;
}

function Metadata({ tabIndex, article, setArticle }: MetadataProps) {
  const [saving, setSaving] = useState<SavingStatus>("untouched");

  return (
    <Flex gap="3" mb="5">
      <Box flexGrow="3">
        <TextField.Root
          tabIndex={tabIndex}
          size="2"
          placeholder="Title"
          value={article.title}
          disabled={saving === "saving"}
          onChange={(e) => {
            setSaving("touched");
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
          disabled={saving === "saving"}
          onChange={(e) => {
            setSaving("touched");
            setArticle({ ...article, author: e.target.value });
          }}
        >
          <TextField.Slot>auteur</TextField.Slot>
        </TextField.Root>
      </Box>
      <Select.Root
        value={article.layout_id + ""}
        onValueChange={(v) => {
          setArticle({ ...article, layout_id: Number(v) });
          setSaving("touched");
        }}
      >
        <Select.Trigger />
        <Select.Content>
          <Select.Group>
            <Select.Label>Layout</Select.Label>
            {article.layouts.map((a) => {
              if (a.value == article.layout_id) {
                return (
                  <Select.Item key={a.value} value={a.value + ""} disabled>
                    {a.name}
                  </Select.Item>
                );
              }
              return (
                <Select.Item key={a.value} value={a.value + ""}>
                  {a.name}
                </Select.Item>
              );
            })}
          </Select.Group>
        </Select.Content>
      </Select.Root>
      <Button
        tabIndex={tabIndex}
        size="2"
        disabled={saving != "touched"}
        onClick={async () => {
          setSaving("saving");
          await postArticleEditMetadata(article.slug, article);
          setSaving("untouched");
        }}
      >
        {saving == "saving" ? <Spinner /> : "Save"}
      </Button>
    </Flex>
  );
}

interface ContentProps {
  tabIndex: number;
  article: Article;
  setArticle: (a: Article) => void;
}

function Content({ tabIndex, article, setArticle }: ContentProps) {
  const [saving, setSaving] = useState<SavingStatus>("untouched");

  return (
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
          style={{ paddingTop: 48 }}
          value={article.content}
          disabled={saving === "saving"}
          onChange={(e) => {
            setSaving("touched");
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
                disabled={saving != "touched"}
                onClick={async () => {
                  setSaving("saving");
                  await postArticleEditContent(
                    article.slug,
                    article.content || "",
                  );
                  setSaving("untouched");
                }}
              >
                {saving == "saving" ? <Spinner /> : "Save"}
              </Button>
            </Flex>
          </Flex>
        </Box>
      </Box>
    </Box>
  );
}

interface AddBlockProps {
  tabIndex: number;
  article: Article;
  setArticle: (a: Article) => void;
}

function AddBlock({ tabIndex, article }: AddBlockProps) {
  const [saving, setSaving] = useState<SavingStatus>("untouched");
  const [block, setBlock] = useState<string>("");

  return (
    <Flex gap="3" mb="5">
      <Select.Root
        value={block}
        onValueChange={(v) => {
          setBlock(v);
          setSaving("touched");
        }}
      >
        <Select.Trigger />
        <Select.Content>
          <Select.Group>
            <Select.Label>Blocks to add</Select.Label>
            {article.blocks.map((b) => {
              if (b.value + "" == block) {
                return (
                  <Select.Item key={b.value} value={b.value + ""} disabled>
                    {b.name}
                  </Select.Item>
                );
              }
              return (
                <Select.Item key={b.value} value={b.value + ""}>
                  {b.name}
                </Select.Item>
              );
            })}
          </Select.Group>
        </Select.Content>
      </Select.Root>
      <Button
        tabIndex={tabIndex}
        size="2"
        disabled={saving != "touched"}
        onClick={async () => {
          setSaving("saving");
          await postArticleEditBlockAdd(article.slug, Number(block));
          setSaving("untouched");
        }}
      >
        {saving == "saving" ? <Spinner /> : "Add"}
      </Button>
    </Flex>
  );
}

export default function Articles() {
  const navigate = useNavigate();
  const { slug } = useParams<{ slug: string }>();
  const [article, setArticle] = useState<Article | null>(null);
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
          <Metadata
            tabIndex={tabIndex}
            article={article}
            setArticle={setArticle}
          ></Metadata>
          <Content
            tabIndex={tabIndex}
            article={article}
            setArticle={setArticle}
          ></Content>
          <AddBlock
            tabIndex={tabIndex}
            article={article}
            setArticle={setArticle}
          ></AddBlock>
          <Flex direction="row" gap="5" wrap="wrap">
            {article.block_datas.map((block) => (
              <DynamicForm
                key={block.id}
                // prefix={`block.${index}`}
                name={block.name}
                data={block.data}
                ui={makeRadixUI(300)}
                // onSubmit={(updated) => {
                //   const newBlockDatas = [...article.block_datas];
                //   newBlockDatas[index] = Object.fromEntries(
                //     Object.entries(updated).map(([k, v]) => [
                //       k.replace(`block.${index}.`, ""),
                //       v,
                //     ]),
                //   );
                //   setArticle({ ...article, block_datas: newBlockDatas });
                // }}
              />
            ))}
          </Flex>

          {/*{article.block_datas.map((b) => {
            return <p>{JSON.stringify(b)}</p>;
          })}*/}
        </Flex>
      </Card>
    </Flex>
  );
}
