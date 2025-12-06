import { Theme, ThemePanel } from "@radix-ui/themes";
import "@radix-ui/themes/styles.css";
import { useTheme } from "./hooks/useTheme";
import App from "./App";

export function Root() {
  const { theme, toggleTheme } = useTheme("light");

  return (
    <Theme appearance={theme} accentColor="indigo" radius="medium">
      <App toggleTheme={toggleTheme} />
      <ThemePanel />
    </Theme>
  );
}
