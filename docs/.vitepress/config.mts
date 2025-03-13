import { defineConfig } from "vitepress";

export default defineConfig({
  base: "/pair/",
  title: "Pair",
  description: "A Git co-authors management tool",
  themeConfig: {
    socialLinks: [
      { icon: "github", link: "https://github.com/philippeckel/pair" },
    ],
    logo: "/favicon.svg",
    logoLink: "/",
    sidebar: [
      {
        text: "Introduction",
        items: [
          { text: "What is Pair?", link: "/about" },
          { text: "Installation", link: "/installation" },
        ],
      },
      {
        text: "Configuration",
        items: [
          { text: "Configuration file", link: "/configuration-file" },
          { text: "Environment variables", link: "/environment-variables" },
        ],
      },
      {
        text: "CLI reference",
        items: [{ text: "pair", link: "/reference/pair" }],
      },
    ],
  },
  head: [["link", { rel: "icon", href: "/favicon.svg" }]],

  markdown: {
    theme: {
      light: "catppuccin-latte",
      dark: "catppuccin-macchiato",
    },
  },
});
