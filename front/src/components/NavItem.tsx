import * as NavigationMenu from "@radix-ui/react-navigation-menu";
import type { ReactNode } from "react";
import { Box, Text } from "@radix-ui/themes";

interface NavItemProps {
  active: boolean;
  onClick: () => void;
  children: ReactNode;
}

export function NavItem({ active, onClick, children }: NavItemProps) {
  return (
    <NavigationMenu.Link onClick={onClick}>
      <Box
        style={{
          cursor: "pointer",
          borderRadius: "var(--radius-3)",
          padding: "var(--space-2)",
          backgroundColor: active ? "var(--accent-a2)" : "transparent",
        }}
      >
        <Text
          size={"4"}
          weight={active ? "bold" : "medium"}
          style={{
            color: active ? "var(--accent-11)" : "var(--gray-12)",
          }}
        >
          {children}
        </Text>
      </Box>
    </NavigationMenu.Link>
  );
}
