import { useEffect, useState } from "react";
import { Text, Flex, Spinner, Card } from "@radix-ui/themes";
import { deleteFile, getFileList, postFile, type File } from "../../api/file";
import Line from "./components/Line";
import SingleFileUploader from "./components/Uploader";

export default function Files() {
  const [files, setFiles] = useState<File[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      try {
        const res = await getFileList();
        setFiles(res.items);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

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
            const res = await getFileList();
            setFiles(res.items);
          }}
        ></SingleFileUploader>
      </Flex>
      <Card>
        <Flex direction="column">
          {files.map((val, idx) => {
            return (
              <Line
                key={idx}
                tabIndex={idx}
                file={val}
                deleteFile={async (filename: string) => {
                  await deleteFile(filename);
                  setFiles(files.filter((f) => f.name != filename));
                }}
                portalContainer={null}
              ></Line>
            );
          })}
        </Flex>
      </Card>
    </Flex>
  );
}
