import { Box, Button, TextArea } from "@radix-ui/themes";
import { TextField } from "@radix-ui/themes";
import { Label } from "@radix-ui/react-label";

import { useEffect, useState } from "react";
import { Text, Flex, Spinner, Card } from "@radix-ui/themes";
import { getBlockEdit, postBlockEdit } from "../../api/block";
import type { Block } from "../../api/block";
import { useNavigate, useParams } from "react-router-dom";
import { Dict } from "../../pkg/data/parseFormData";

type SavingStatus = "untouched" | "touched" | "saving";

interface NameProps {
  tabIndex: number;
  slug: string;
  block: Block;
  setBlock: (b: Block) => void;
}

function Name({ tabIndex, slug, block, setBlock }: NameProps) {
  const [saving, setSaving] = useState<SavingStatus>("untouched");
  const navigate = useNavigate();
  return (
    <Flex gap="3" mb="5">
      <Box flexGrow="3">
        <TextField.Root
          tabIndex={tabIndex}
          size="2"
          placeholder="Title"
          value={block.name}
          disabled={saving === "saving"}
          onChange={(e) => {
            setSaving("touched");
            setBlock({ ...block, name: e.target.value });
          }}
        >
          <TextField.Slot>name</TextField.Slot>
        </TextField.Root>
      </Box>
      <Button
        tabIndex={tabIndex}
        size="2"
        disabled={saving != "touched"}
        onClick={async () => {
          setSaving("saving");
          await postBlockEdit(slug, block);
          setSaving("untouched");
          navigate(`/blocks/${block.name}/edit`, { replace: true });
        }}
      >
        {saving == "saving" ? <Spinner /> : "Save"}
      </Button>
    </Flex>
  );
}

interface EditorProps {
  tabIndex: number;
  title: string;
  content: string;
  setContent: (content: string) => void;
  updateContent: () => Promise<void>;
}

function Editor({
  tabIndex,
  title,
  content,
  setContent,
  updateContent,
}: EditorProps) {
  const [saving, setSaving] = useState<SavingStatus>("untouched");
  return (
    <Box mb="4">
      <Label htmlFor="skirt-description">
        <Text size="2" weight="bold" mb="2" asChild>
          <Box display="inline-block">{title}</Box>
        </Text>
      </Label>
      <Box position="relative">
        <TextArea
          tabIndex={tabIndex}
          spellCheck={false}
          id="skirt-description"
          variant="soft"
          rows={10}
          value={content}
          style={{ paddingTop: 48 }}
          disabled={saving === "saving"}
          onChange={(e) => {
            setSaving("touched");
            setContent(e.target.value);
          }}
        />
        <Box position="absolute" m="2" top="0" left="0" right="0">
          <Flex gap={"1"}>
            <Button
              tabIndex={tabIndex}
              size="2"
              disabled={saving != "touched"}
              onClick={async () => {
                setSaving("saving");
                await updateContent();
                setSaving("untouched");
              }}
            >
              {saving == "saving" ? <Spinner /> : "Save"}
            </Button>
          </Flex>
        </Box>
      </Box>
    </Box>
  );
}

export default function BlockEdit() {
  const { slug } = useParams<{ slug: string }>();
  const [block, setBlock] = useState<Block | null>(null);
  const [temp, setTemp] = useState<string | null>(null);
  const navidate = useNavigate();
  useEffect(() => {
    async function load() {
      if (!slug) {
        setBlock(null);
        return;
      }
      const res = await getBlockEdit(slug);
      setBlock(res);
    }
    load();
  }, [slug]);

  if (!slug) {
    navidate("/", { replace: true });
    return (
      <Flex align="center" justify="center" height="100vh">
        <Spinner />
      </Flex>
    );
  }

  if (!block) {
    return (
      <Flex align="center" justify="center" height="100vh">
        <Spinner />
      </Flex>
    );
  }

  return (
    <Flex direction="column" gap="4">
      <Text size="7" weight="bold">
        {slug}
      </Text>
      <Card>
        <Flex direction="column">
          <Name
            slug={slug}
            tabIndex={1}
            block={block}
            setBlock={setBlock}
          ></Name>
          <Editor
            tabIndex={2}
            title="Content"
            content={block.content}
            setContent={(content: string) => {
              setBlock({ ...block, content: content });
            }}
            updateContent={async () => {
              await postBlockEdit(slug, block);
            }}
          ></Editor>
          <Editor
            tabIndex={3}
            title="Data"
            content={temp || JSON.stringify(block.definition, null, 2)}
            setContent={(content: string) => {
              let p: Dict | null = null;
              try {
                p = JSON.parse(content);
              } catch {
                setTemp(content);
              } finally {
                if (p != null) {
                  setTemp(null);
                  setBlock({ ...block, definition: p });
                }
              }
            }}
            updateContent={async () => {
              await postBlockEdit(slug, block);
            }}
          ></Editor>
        </Flex>
      </Card>
    </Flex>
  );
}
