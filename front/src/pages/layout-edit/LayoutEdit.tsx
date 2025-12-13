import { Box, Button, TextArea } from "@radix-ui/themes";
import { TextField } from "@radix-ui/themes";
import { Label } from "@radix-ui/react-label";

import { useEffect, useState } from "react";
import { Text, Flex, Spinner, Card } from "@radix-ui/themes";
import { getLayoutEdit, postLayoutEdit } from "../../api/layout";
import type { Layout } from "../../api/layout";
import { useParams } from "react-router-dom";

type SavingStatus = "untouched" | "touched" | "saving";

export default function LayoutEdit() {
  const { slug } = useParams<{ slug: string }>();
  const [layoutName, setLayoutName] = useState<string>("");
  const [layout, setLayout] = useState<Layout | null>(null);
  const [saving, setSaving] = useState<SavingStatus>("untouched");
  const tabIndex = 1;
  useEffect(() => {
    async function load() {
      if (saving !== "untouched") {
        return;
      }
      if (!slug) {
        setLayout(null);
        setLayoutName("");
        return;
      }
      const res = await getLayoutEdit(slug);
      setLayout(res);
      setLayoutName(res.name);
    }
    load();
  }, [slug, saving]);

  if (!layout) {
    return (
      <Flex align="center" justify="center" height="100vh">
        <Spinner />
      </Flex>
    );
  }

  return (
    <Flex direction="column" gap="4">
      <Text size="7" weight="bold">
        {layoutName}
      </Text>
      <Card>
        <Flex direction="column">
          <Flex gap="3" mb="5">
            <Box flexGrow="2">
              <TextField.Root
                tabIndex={tabIndex}
                size="2"
                placeholder="Title"
                value={layout.name}
                disabled={saving === "saving"}
                onChange={(e) => {
                  setSaving("touched");
                  setLayout({ ...layout, name: e.target.value });
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
                await postLayoutEdit(layoutName, layout);
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
            <TextArea
              tabIndex={tabIndex}
              spellCheck={false}
              id="skirt-description"
              variant="soft"
              rows={25}
              defaultValue={layout.content || ""}
              value={layout.content}
              disabled={saving === "saving"}
              onChange={(e) => {
                setSaving("touched");
                setLayout({ ...layout, content: e.target.value });
              }}
            />
          </Box>
        </Flex>
      </Card>
    </Flex>
  );
}
