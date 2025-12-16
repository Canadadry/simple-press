import { useRef, useState } from "react";
import {
  Button,
  Flex,
  Text,
  Card,
  TextField,
  Checkbox,
} from "@radix-ui/themes";

interface SingleFileUploaderProps {
  handleUpload: (f: Blob, fileName: string, archive: boolean) => Promise<void>;
}

function SingleFileUploader({ handleUpload }: SingleFileUploaderProps) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [file, setFile] = useState<File | null>(null);
  const [archive, setArchive] = useState<boolean | null>(null);
  const [fileName, setFileName] = useState<string>("");

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0] ?? null;
    setFile(selectedFile);
    setFileName(selectedFile?.name ?? "");
    setArchive(selectedFile?.name.endsWith(".zip") ? false : null);
  };

  const handleSubmit = async () => {
    if (!file) return;
    await handleUpload(file, fileName, archive == true);
    setFile(null);
    setFileName("");
    if (inputRef.current) {
      inputRef.current.value = "";
    }
  };

  return (
    <Flex direction="column" gap="3">
      {/* Input file natif */}
      <input ref={inputRef} type="file" hidden onChange={handleFileChange} />

      <Button variant="soft" onClick={() => inputRef.current?.click()}>
        Choisir un fichier
      </Button>

      {file && (
        <Card>
          <Flex direction="column" gap="2">
            <Text weight="bold" size="2">
              Détails du fichier
            </Text>

            {/* Champ éditable */}
            <TextField.Root
              value={fileName}
              onChange={(e) => setFileName(e.target.value)}
              placeholder="Nom du fichier"
            />

            <Text size="2">Type : {file.type || "—"}</Text>
            <Text size="2">Taille : {file.size} bytes</Text>
            {archive != null ? (
              <Text as="label" size="2" data-testid={`checkbox-${name}`}>
                <Flex gap="2">
                  <Checkbox
                    mb="2"
                    checked={archive}
                    onCheckedChange={(c) => {
                      setArchive(c === "indeterminate" ? false : c);
                    }}
                  />
                  {"explode archive after upload"}
                </Flex>
              </Text>
            ) : (
              <></>
            )}
          </Flex>
        </Card>
      )}

      {file && <Button onClick={handleSubmit}>Upload le fichier</Button>}
    </Flex>
  );
}

export default SingleFileUploader;
