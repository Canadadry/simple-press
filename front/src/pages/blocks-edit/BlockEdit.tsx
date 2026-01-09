import { Box, Button } from "@radix-ui/themes";
import { TextField } from "@radix-ui/themes";

import { useEffect, useState } from "react";
import { Text, Flex, Spinner, Card } from "@radix-ui/themes";
import { getBlockEdit, postBlockEdit } from "../../api/block";
import type { Block } from "../../api/block";
import { useNavigate, useParams } from "react-router-dom";
import { Editor } from "../../components/Editor";
import { JsonEditor } from "../../components/JsonEditor";
import { Dict } from "../../api/api";

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

export default function BlockEdit() {
  const { slug } = useParams<{ slug: string }>();
  const [block, setBlock] = useState<Block | null>(null);
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
          <JsonEditor
            tabIndex={3}
            title="Data"
            content={block.definition}
            setContent={(content: Dict) => {
              setBlock({ ...block, definition: content });
            }}
            updateContent={async () => {
              await postBlockEdit(slug, block);
            }}
          ></JsonEditor>
        </Flex>
      </Card>
    </Flex>
  );
}
