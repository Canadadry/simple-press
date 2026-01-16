import { DotsHorizontalIcon } from "@radix-ui/react-icons";
import { PersonIcon } from "@radix-ui/react-icons";
import {
  Avatar,
  Box,
  DropdownMenu,
  IconButton,
  Separator,
  Text,
  Flex,
} from "@radix-ui/themes";
import type { Article } from "../../../api/article";
import { useNavigate } from "react-router-dom";
import { useState } from "react";

interface LineProps {
  index: number;
  article: Article;
  portalContainer?: Element | DocumentFragment | null;
}
export default function Line(line: LineProps) {
  const navigate = useNavigate();
  const [isHovered, setIsHovered] = useState<boolean>(false);

  return (
    <Box key={line.index}>
      {line.index && line.index > 0 ? (
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
          navigate(`/articles/${line.article.slug}/edit`, {
            replace: true,
          });
        }}
      >
        <Flex gap="3" align="center" width="200px">
          <Avatar
            src={line.article.image}
            fallback={line.article.title[0].toUpperCase()}
          />
          <Text
            size="2"
            weight={"bold"}
            style={{
              color: "var(--accent-11)",
            }}
          >
            {line.article.title}
          </Text>
        </Flex>
        <PersonIcon></PersonIcon>
        <Text size="2" color="gray">
          {line.article.author}
        </Text>
        <Text size="2">{line.article.content.slice(0, 50)}...</Text>

        <Flex flexGrow="1" justify="end">
          <DropdownMenu.Root>
            <DropdownMenu.Trigger>
              <IconButton
                color="gray"
                tabIndex={line.index + 3}
                variant="ghost"
              >
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
