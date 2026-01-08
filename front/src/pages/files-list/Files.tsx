import { useEffect, useState } from "react";
import { Text, Flex, Spinner, Card } from "@radix-ui/themes";
import {
  deleteFile,
  type FileTree,
  getFileList,
  postFile,
  type File,
  getFileTree,
} from "../../api/file";
import Line from "./components/Line";
import SingleFileUploader from "./components/Uploader";
import { Button } from "@radix-ui/themes/dist/cjs/index.js";

function removeLastPathSegment(path: string): string {
  if (!path) {
    return path;
  }

  const normalizedPath = path.replace(/\/+$/, "");
  const lastSlashIndex = normalizedPath.lastIndexOf("/");

  if (lastSlashIndex === -1) {
    return "";
  }

  return normalizedPath.substring(0, lastSlashIndex);
}

export default function Files() {
  const [loading, setLoading] = useState<boolean>(true);
  const [path, setPath] = useState<string>("");
  const [tree, setTree] = useState<FileTree>({
    path: "",
    folders: [],
    files: [],
  });

  useEffect(() => {
    async function load() {
      try {
        const resTree = await getFileTree(path);
        setTree(resTree);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, [path]);

  if (loading) {
    return (
      <Flex align="center" justify="center" height="100vh">
        <Spinner />
      </Flex>
    );
  }

  return (
    <Flex direction="column" gap="4">
      <Text size="7" weight="bold">
        Liste des files
      </Text>
      <Flex maxWidth={"500"}>
        <SingleFileUploader
          handleUpload={async (f: Blob, filename: string, achive: boolean) => {
            await postFile(f, filename, achive);
          }}
        ></SingleFileUploader>
      </Flex>
      <Card>
        <p>{path}</p>
        {path === "" ? (
          <></>
        ) : (
          <Button
            onClick={() => {
              setPath(removeLastPathSegment(path));
            }}
          >
            ..
          </Button>
        )}
        {tree.folders.map((f: string) => {
          return (
            <Button
              onClick={() => {
                setPath(path + "/" + f);
              }}
            >
              {f}
            </Button>
          );
        })}
      </Card>
      <Card>
        {
          <Flex direction="column">
            {tree.files.map((val, idx) => {
              return (
                <Line
                  key={idx}
                  tabIndex={idx}
                  file={val}
                  path={path}
                  deleteFile={async (filename: string) => {
                    await deleteFile(filename);
                    setTree({
                      ...tree,
                      files: tree.files.filter((f) => f != filename),
                    });
                  }}
                  portalContainer={null}
                ></Line>
              );
            })}
          </Flex>
        }
      </Card>
    </Flex>
  );
}
