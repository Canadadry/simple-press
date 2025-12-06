import { useState, useEffect } from "react";

export type ThemeMode = "light" | "dark";

export function useTheme(defaultTheme: ThemeMode = "light") {
  const [theme, setTheme] = useState<ThemeMode>(() => {
    // lecture du localStorage au premier render
    const saved = localStorage.getItem("theme") as ThemeMode;
    return saved ?? defaultTheme;
  });

  // sauvegarder à chaque changement
  useEffect(() => {
    localStorage.setItem("theme", theme);
  }, [theme]);

  const toggleTheme = () => {
    setTheme((prev) => (prev === "light" ? "dark" : "light"));
  };

  return { theme, toggleTheme };
}
