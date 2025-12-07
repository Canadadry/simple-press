import { DotsHorizontalIcon } from "@radix-ui/react-icons";
import {
  Avatar,
  Box,
  DropdownMenu,
  IconButton,
  Link,
  Separator,
} from "@radix-ui/themes";
import { Text, Flex } from "@radix-ui/themes";
import type { Article } from "../../../api/article";
import { useNavigate } from "react-router-dom";

interface LineProps {
  tabIndex: number | undefined;
  article: Article;
  portalContainer: Element | DocumentFragment | null | undefined;
}
export default function Line(line: LineProps) {
  const navigate = useNavigate();

  return (
    <Box key={line.tabIndex}>
      {line.tabIndex && line.tabIndex > 0 ? (
        <Box>
          <Separator size="4" my="3" />
        </Box>
      ) : (
        <></>
      )}
      <Flex gap="4" align="center">
        <Flex gap="3" align="center" width="200px">
          <Avatar
            src={line.article.image}
            fallback={line.article.title[0].toUpperCase()}
          />
          <Link
            href="#"
            target="_blank"
            tabIndex={line.tabIndex}
            size="2"
            wrap="nowrap"
            onClick={(e) => {
              e.preventDefault();
              navigate(`/articles/${line.article.slug}/edit`, {
                replace: true,
              });
            }}
          >
            {line.article.title}
          </Link>
        </Flex>

        <Text size="2" color="gray">
          {line.article.author}
        </Text>

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
