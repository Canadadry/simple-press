import { Text, Flex } from "@radix-ui/themes";
import { DynamicForm } from "../../pkg/data/render";
import { useEffect, useState } from "react";
import { Dict } from "../../api/api";
import { makeRadixUI } from "../../pkg/data/radix-form";
import { getGlobalData } from "../../api/global";

export default function GlobalData() {
  const [data, setData] = useState<Dict>({});
  useEffect(() => {
    const func = async () => {
      const d = await getGlobalData();
      setData(d);
    };
    func();
  }, []);

  return (
    <Flex direction="column" gap="4">
      <Text size="7" weight="bold">
        Web Site data
      </Text>
      <DynamicForm
        name="GlobalData"
        data={data}
        ui={makeRadixUI(500)}
        setData={(d) => {
          setData(d);
        }}
        onSave={async () => {}}
      />
    </Flex>
  );
}
