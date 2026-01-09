import { Box, Button, TextArea } from "@radix-ui/themes";
import { Label } from "@radix-ui/react-label";

import { useState } from "react";
import { Text, Flex, Spinner } from "@radix-ui/themes";

type SavingStatus = "untouched" | "touched" | "saving";

export interface EditorProps {
  tabIndex: number;
  title: string;
  content: string;
  setContent: (content: string) => void;
  updateContent: () => Promise<void>;
}

export const Editor: React.FC<EditorProps> = ({
  tabIndex,
  title,
  content,
  setContent,
  updateContent,
}: EditorProps) => {
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
};
