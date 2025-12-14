import { useEffect, useState } from "react";
import {
  Text,
  Flex,
  Spinner,
  Card,
  Box,
  Button,
  TextField,
} from "@radix-ui/themes";
// import * as Accordion from "@radix-ui/react-accordion";
import { getBlockList, postBlockAdd, type Block } from "../../api/block";
import Line from "./components/Line";
import { useNavigate } from "react-router-dom";

type SavingStatus = "untouched" | "touched" | "saving";

interface CreateProps {
  tabIndex: number;
}

function Create({ tabIndex }: CreateProps) {
  const [saving, setSaving] = useState<SavingStatus>("untouched");
  const [name, setName] = useState<string>("");
  const navigate = useNavigate();
  return (
    <Flex gap="3" mb="5">
      <Box flexGrow="3">
        <TextField.Root
          tabIndex={tabIndex}
          size="2"
          placeholder="Title"
          value={name}
          disabled={saving === "saving"}
          onChange={(e) => {
            setSaving("touched");
            setName(e.target.value);
          }}
        >
          <TextField.Slot>titre</TextField.Slot>
        </TextField.Root>
      </Box>
      <Button
        tabIndex={tabIndex}
        size="2"
        disabled={saving != "touched"}
        onClick={async () => {
          setSaving("saving");
          await postBlockAdd(name);
          navigate(`/blocks/${name}/edit`, { replace: true });
          setSaving("untouched");
        }}
      >
        {saving == "saving" ? <Spinner /> : "Create"}
      </Button>
    </Flex>
  );
}

export default function Blocks() {
  const [blocks, setBlocks] = useState<Block[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      try {
        const res = await getBlockList();
        setBlocks(res.items);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  if (loading) {
    return (
      <Flex direction="column" gap="4">
        <Text size="7" weight="bold">
          Liste des blocks
        </Text>
        <Create tabIndex={0}></Create>
        <Flex align="center" justify="center" height="100vh">
          <Spinner />
        </Flex>
      </Flex>
    );
  }

  return (
    <Flex direction="column" gap="4">
      <Text size="7" weight="bold">
        Liste des blocks
      </Text>
      <Create tabIndex={0}></Create>
      <Card>
        <Flex direction="column">
          {blocks.map((val, idx) => {
            return (
              <Line
                key={idx}
                tabIndex={idx}
                block={val}
                portalContainer={null}
              ></Line>
            );
          })}
        </Flex>
      </Card>
    </Flex>
  );
}
