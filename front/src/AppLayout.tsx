import { Flex, Box, Text } from "@radix-ui/themes";
import * as NavigationMenu from "@radix-ui/react-navigation-menu";
import { useNavigate, useLocation } from "react-router-dom";
import React from "react";
import { Button } from "@radix-ui/themes";

interface AppLayoutProps {
  children: React.ReactNode;
  toggleTheme: () => void;
}

export function AppLayout({ children, toggleTheme }: AppLayoutProps) {
  const navigate = useNavigate();
  const location = useLocation();

  const navLinkStyle = (path: string): React.CSSProperties => ({
    padding: "8px 0",
    display: "block",
    fontSize: "16px",
    textDecoration: "none",
    cursor: "pointer",
    color: location.pathname === path ? "var(--indigo-11)" : "var(--gray-12)",
    fontWeight: location.pathname === path ? 600 : 400,
  });

  return (
    <Flex width="100%" height="100vh">
      {/* Sidebar */}
      <Box
        width="240px"
        p="4"
        style={{
          borderRight: "1px solid var(--gray-a5)",
          boxSizing: "border-box",
        }}
      >
        <Text size="5" weight="bold" mb="4" as="div">
          Mon SaaS
        </Text>

        <NavigationMenu.Root orientation="vertical">
          <NavigationMenu.List style={{ listStyle: "none", padding: 0 }}>
            <NavigationMenu.Item>
              <NavigationMenu.Link
                style={navLinkStyle("/dashboard")}
                onClick={() => navigate("/dashboard")}
              >
                Dashboard
              </NavigationMenu.Link>
            </NavigationMenu.Item>
            <NavigationMenu.Item>
              <NavigationMenu.Link
                style={navLinkStyle("/users")}
                onClick={() => navigate("/users")}
              >
                Utilisateurs
              </NavigationMenu.Link>
            </NavigationMenu.Item>
            <NavigationMenu.Item>
              <NavigationMenu.Link
                style={navLinkStyle("/settings")}
                onClick={() => navigate("/settings")}
              >
                Paramètres
              </NavigationMenu.Link>
            </NavigationMenu.Item>
          </NavigationMenu.List>
        </NavigationMenu.Root>
        <Box mt="auto">
          <Button onClick={toggleTheme} size="2">
            Changer de thème
          </Button>
        </Box>
      </Box>
      {/* Main content */}
      <Box flexGrow="1" p="6">
        {children}
      </Box>
    </Flex>
  );
}
