import { DotsHorizontalIcon, Pencil1Icon } from "@radix-ui/react-icons";
import {
  Avatar,
  Box,
  DropdownMenu,
  IconButton,
  Separator,
  Text,
  Flex,
} from "@radix-ui/themes";

interface LineProps {
  path: string;
  tabIndex: number | undefined;
  file: string;
  portalContainer: Element | DocumentFragment | null | undefined;
  deleteFile: (filename: string) => Promise<void>;
}
export default function Line(line: LineProps) {
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
        <Flex gap="3" align="center">
          <Avatar src={line.file} fallback={line.file[0].toUpperCase()} />
          <Text
            size="2"
            weight={"bold"}
            style={{
              color: "var(--accent-11)",
            }}
          >
            {line.file}
          </Text>
          <Pencil1Icon></Pencil1Icon>
        </Flex>
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
              {/*<DropdownMenu.Item>Preview</DropdownMenu.Item>
              <DropdownMenu.Item>Edit</DropdownMenu.Item>
              <DropdownMenu.Separator />*/}
              <DropdownMenu.Item
                color="red"
                onClick={(e) => {
                  e.stopPropagation();
                  line.deleteFile(line.path + "/" + line.file);
                }}
              >
                Remove
              </DropdownMenu.Item>
            </DropdownMenu.Content>
          </DropdownMenu.Root>
        </Flex>
      </Flex>
    </Box>
  );
}
