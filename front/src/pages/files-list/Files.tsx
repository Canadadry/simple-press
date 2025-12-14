import { useEffect, useState } from "react";
import { Text, Flex, Spinner, Card } from "@radix-ui/themes";
// import * as Accordion from "@radix-ui/react-accordion";
import { getFileList, postFile, type File } from "../../api/file";
import Line from "./components/Line";
import SingleFileUploader from "./components/Uploader";
import { useNavigate } from "react-router-dom";

export default function Files() {
  const [files, setFiles] = useState<File[]>([]);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

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
      <SingleFileUploader
        handleUpload={async (f: Blob) => {
          postFile(f);
          const res = await getFileList();
          setFiles(res.items);
        }}
      ></SingleFileUploader>
      <Card>
        <Flex direction="column">
          {files.map((val, idx) => {
            return (
              <Line
                key={idx}
                tabIndex={idx}
                file={val}
                portalContainer={null}
              ></Line>
            );
          })}
        </Flex>
      </Card>
    </Flex>
  );
}
