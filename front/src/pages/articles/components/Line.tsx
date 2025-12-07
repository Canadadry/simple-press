import { DotsHorizontalIcon } from "@radix-ui/react-icons";
import {
  Avatar,
  Box,
  Button,
  DropdownMenu,
  Heading,
  IconButton,
  Link,
  Separator,
  TextField,
} from "@radix-ui/themes";
import { Text, Flex, Card } from "@radix-ui/themes";
import type { Article } from "../../../api/article";

interface LineProps {
  tabIndex: number | undefined;
  article: Article;
  portalContainer: Element | DocumentFragment | null | undefined;
}
export default function Line(line: LineProps) {
  return (
    <Card size="4">
      <Heading as="h3" size="6" trim="start" mb="2">
        Your team
      </Heading>

      <Text as="p" size="2" mb="5" color="gray">
        Invite and manage your team members.
      </Text>

      <Flex gap="3" mb="5">
        <Box flexGrow="1">
          <TextField.Root
            tabIndex={line.tabIndex}
            size="2"
            placeholder="Email address"
          />
        </Box>
        <Button tabIndex={line.tabIndex} size="2">
          Invite
        </Button>
      </Flex>

      <Flex direction="column">
        {[4, 2, 12, 20, 16].map((number, i) => (
          <Box key={number}>
            <Flex gap="4" align="center">
              <Flex gap="3" align="center" width="200px">
                <Avatar
                  src={line.article.image}
                  fallback={line.article.title[0].toUpperCase()}
                />
                <Link
                  href={
                    import.meta.env.VITE_API_URL +
                    "/admin/articles/" +
                    line.article.slug +
                    "/preview"
                  }
                  target="_blank"
                  tabIndex={line.tabIndex}
                  size="2"
                  wrap="nowrap"
                  onClick={(e) => e.preventDefault()}
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
                    <IconButton
                      color="gray"
                      tabIndex={line.tabIndex}
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

            {i !== 4 && (
              <Box>
                <Separator size="4" my="3" />
              </Box>
            )}
          </Box>
        ))}
      </Flex>
    </Card>
  );
}
