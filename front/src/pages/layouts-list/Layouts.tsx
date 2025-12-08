import { useEffect, useState } from "react";
import { Text, Flex, Spinner, Card } from "@radix-ui/themes";
// import * as Accordion from "@radix-ui/react-accordion";
import { getLayoutList, type Layout } from "../../api/layout";
import Line from "./components/Line";

export default function Layouts() {
  const [layouts, setLayouts] = useState<Layout[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function load() {
      try {
        const res = await getLayoutList();
        setLayouts(res.items);
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
        Liste des layouts
      </Text>
      <Card>
        <Flex direction="column">
          {layouts.map((val, idx) => {
            return (
              <Line
                key={idx}
                tabIndex={idx}
                layout={val}
                portalContainer={null}
              ></Line>
            );
          })}
        </Flex>
      </Card>
    </Flex>
  );
}
