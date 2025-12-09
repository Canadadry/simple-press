import { DotsHorizontalIcon } from "@radix-ui/react-icons";
import {
  Box,
  DropdownMenu,
  IconButton,
  Separator,
  Text,
  Flex,
} from "@radix-ui/themes";
import type { Layout } from "../../../api/layout";
import { useNavigate } from "react-router-dom";
import { useState } from "react";

interface LineProps {
  tabIndex: number | undefined;
  layout: Layout;
  portalContainer: Element | DocumentFragment | null | undefined;
}
export default function Line(line: LineProps) {
  const navigate = useNavigate();
  const [isHovered, setIsHovered] = useState<boolean>(false);

  return (
    <Box key={line.tabIndex}>
      {line.tabIndex && line.tabIndex > 0 ? (
        <Box>
          <Separator size="4" my="3" />
        </Box>
      ) : (
        <></>
      )}
      <Flex
        gap="4"
        align="center"
        style={{
          cursor: "pointer",
          backgroundColor: isHovered ? "var(--accent-a2)" : "transparent",
        }}
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
        onClick={(e) => {
          e.preventDefault();
          navigate(`/layouts/${line.layout.id}/edit`, {
            replace: true,
          });
        }}
      >
        <Flex gap="3" align="center" width="200px">
          <Text
            size="2"
            weight={"bold"}
            style={{
              color: "var(--accent-11)",
            }}
          >
            {line.layout.name}
          </Text>
        </Flex>
        <Text size="2">{line.layout.content.slice(0, 50)}...</Text>

        <Flex flexGrow="1" justify="end">
          <DropdownMenu.Root>
            <DropdownMenu.Trigger>
              <IconButton color="gray" tabIndex={line.tabIndex} variant="ghost">
                <DotsHorizontalIcon />
              </IconButton>
            </DropdownMenu.Trigger>
            <DropdownMenu.Content
              container={line.portalContainer}
              variant="soft"
            >
              <DropdownMenu.Item>Preview</DropdownMenu.Item>
              <DropdownMenu.Item>Edit</DropdownMenu.Item>
              <DropdownMenu.Separator />
              <DropdownMenu.Item color="red">Remove</DropdownMenu.Item>
            </DropdownMenu.Content>
          </DropdownMenu.Root>
        </Flex>
      </Flex>
    </Box>
  );
}
