import { Box, Button, TextArea } from "@radix-ui/themes";
import { TextField } from "@radix-ui/themes";
import { Label } from "@radix-ui/react-label";

import { useEffect, useState } from "react";
import { Text, Flex, Spinner, Card } from "@radix-ui/themes";
import { getBlockEdit, postBlockEdit } from "../../api/block";
import type { Block } from "../../api/block";
import { useParams } from "react-router-dom";

type SavingStatus = "untouched" | "touched" | "saving";

export default function BlockEdit() {
  const { slug } = useParams<{ slug: string }>();
  const [blockName, setBlockName] = useState<string>("");
  const [block, setBlock] = useState<Block | null>(null);
  const [saving, setSaving] = useState<SavingStatus>("untouched");
  const tabIndex = 1;
  useEffect(() => {
    async function load() {
      if (saving !== "untouched") {
        return;
      }
      if (!slug) {
        setBlock(null);
        setBlockName("");
        return;
      }
      const res = await getBlockEdit(slug);
      setBlock(res);
      setBlockName(res.name);
    }
    load();
  }, [slug, saving]);

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
        {blockName}
      </Text>
      <Card>
        <Flex direction="column">
          <Flex gap="3" mb="5">
            <Box flexGrow="2">
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
                <TextField.Slot>Nom</TextField.Slot>
              </TextField.Root>
            </Box>
            <Button
              tabIndex={tabIndex}
              size="2"
              disabled={saving != "touched"}
              onClick={async () => {
                setSaving("saving");
                await postBlockEdit(blockName, block);
                setSaving("untouched");
              }}
            >
              {saving == "saving" ? <Spinner /> : "Save"}
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
                defaultValue={block.content || ""}
                value={block.content}
                style={{ paddingTop: 48 }}
                disabled={saving === "saving"}
                onChange={(e) => {
                  setSaving("touched");
                  setBlock({ ...block, content: e.target.value });
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
                      await postBlockEdit(block.name, block);
                      setSaving("untouched");
                    }}
                  >
                    {saving == "saving" ? <Spinner /> : "Save"}
                  </Button>
                </Flex>
              </Box>
            </Box>
          </Box>
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
                defaultValue={JSON.stringify(block.definition) || ""}
                value={JSON.stringify(block.definition)}
                style={{ paddingTop: 48 }}
                disabled={saving === "saving"}
                onChange={(e) => {
                  setSaving("touched");
                  setBlock({
                    ...block,
                    definition: JSON.parse(e.target.value),
                  });
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
                      await postBlockEdit(block.name, block);
                      setSaving("untouched");
                    }}
                  >
                    {saving == "saving" ? <Spinner /> : "Save"}
                  </Button>
                </Flex>
              </Box>
            </Box>
          </Box>
        </Flex>
      </Card>
    </Flex>
  );
}
