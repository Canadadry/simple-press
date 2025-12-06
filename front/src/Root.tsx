import { Theme } from "@radix-ui/themes";
// import { ThemePanel } from "@radix-ui/themes";
import "@radix-ui/themes/styles.css";
import { useTheme } from "./hooks/useTheme";
import App from "./App";

export function Root() {
  const { theme, toggleTheme } = useTheme("dark");

  return (
    <Theme
      appearance={theme}
      accentColor="purple"
      grayColor="mauve"
      radius="large"
    >
      <App toggleTheme={toggleTheme} />
      {/*<ThemePanel />*/}
    </Theme>
  );
}
