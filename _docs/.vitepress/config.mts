import { defineConfig } from "vitepress";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Calendar",
  base: "/calendar/",
  titleTemplate: false,
  cleanUrls: true,
  description: "Calendar Service",
  markdown: {
    theme: "material-theme",
  },
  head: [["link", { rel: "icon", href: "/favicon.ico", type: "image/x-icon" }]],
  themeConfig: {
    search: {
      provider: "local",
    },
    nav: [{ text: "About", link: "/" }],
    sidebar: [{ text: "Quickstart", link: "/quickstart" }],

    socialLinks: [
      { icon: "github", link: "https://github.com/worldline-go/calendar" },
    ],

    editLink: {
      pattern: "https://github.com/worldline-go/calendar/edit/main/_docs/:path",
    },

    lastUpdated: {
      text: "Updated at",
      formatOptions: {
        dateStyle: "full",
        timeStyle: "short",
      },
    },
  },
});
