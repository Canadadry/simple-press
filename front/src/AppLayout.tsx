import { Flex, Box, Text, Button } from "@radix-ui/themes";
import { useNavigate, useLocation } from "react-router-dom";
import { NavItem } from "./components/NavItem";
import * as NavigationMenu from "@radix-ui/react-navigation-menu";

interface AppLayoutProps {
  children: React.ReactNode;
  toggleTheme: () => void;
}

export function AppLayout({ children, toggleTheme }: AppLayoutProps) {
  const navigate = useNavigate();
  const location = useLocation();

  const menu = [
    { label: "Dashboard", path: "/dashboard" },
    { label: "Articles", path: "/articles" },
    { label: "Layouts", path: "/layouts" },
  ];

  return (
    <Flex width="100%" height="100vh">
      {/* Sidebar */}
      <Flex
        direction="column"
        width="240px"
        p="4"
        style={{ borderRight: "1px solid var(--gray-a5)" }}
      >
        <Flex direction="row" gap="2" mb="4" align="center">
          <img
            src="/logo.svg"
            alt="Logo Mon SaaS"
            style={{ width: 36, height: 36 }}
          />
          <Text size="5" weight="bold">
            Mon SaaS
          </Text>
        </Flex>

        <NavigationMenu.Root orientation="vertical">
          <NavigationMenu.List style={{ listStyle: "none", padding: 0 }}>
            {menu.map((item) => (
              <NavigationMenu.Item key={item.path}>
                <NavItem
                  active={location.pathname === item.path}
                  onClick={() => navigate(item.path)}
                >
                  {item.label}
                </NavItem>
              </NavigationMenu.Item>
            ))}
          </NavigationMenu.List>
        </NavigationMenu.Root>

        {/* Theme switcher en bas */}
        <Box mt="auto">
          <Button onClick={toggleTheme} size="2">
            Changer de thème
          </Button>
        </Box>
      </Flex>

      {/* Main content */}
      <Box flexGrow="1" p="6">
        {children}
      </Box>
    </Flex>
  );
}
