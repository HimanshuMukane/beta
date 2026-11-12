import { useEffect, useState } from "react"
import { IconMoon, IconSun } from "../icons/tabler-icons"
import { Button } from "./ui/button"

export function ModeToggle() {
  const [theme, setTheme] = useState<"light" | "dark">(() => {
    if (typeof window !== "undefined") {
      const saved = localStorage.getItem("theme")
      if (saved === "dark" || saved === "light") return saved
      if (window.matchMedia("(prefers-color-scheme: dark)").matches) return "dark"
    }
    return "light"
  })

  useEffect(() => {
    const root = document.documentElement
    if (theme === "dark") {
      root.classList.add("dark")
    } else {
      root.classList.remove("dark")
    }
    localStorage.setItem("theme", theme)
  }, [theme])

  const toggleTheme = () => {
    setTheme((prev) => (prev === "dark" ? "light" : "dark"))
  }

  return (
    <Button
      variant="ghost"
      size="icon"
      onClick={toggleTheme}
      className="h-8 w-8 rounded-md"
      title={theme === "dark" ? "Switch to light mode" : "Switch to dark mode"}
    >
      {theme === "dark" ? (
        <IconSun className="h-4 w-4 transition-all" />
      ) : (
        <IconMoon className="h-4 w-4 transition-all" />
      )}
      <span className="sr-only">Toggle theme</span>
    </Button>
  )
}

