import { Text, Card, Flex, Button } from "@radix-ui/themes";

export default function Dashboard() {
  return (
    <Flex direction="column" gap="4">
      <Text size="7" weight="bold">
        Dashboard
      </Text>
      <Card size="3">
        <Flex direction="column" gap="3">
          <Text size="4">Bienvenue sur le dashboard !</Text>
          <Button>Action du widget</Button>
        </Flex>
      </Card>
    </Flex>
  );
}
