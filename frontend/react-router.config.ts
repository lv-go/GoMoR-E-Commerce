import type { Config } from "@react-router/dev/config";

export default {
  // Config options...
  // Server-side render by default, to enable SPA mode set this to `false`
  ssr: false,
  prerender: [
    "/",
    "/shop",
    "/product/:id",
    "/cart",
    "/favorites",
    "/profile",
    "/user/orders",
    "/login",
    "/register",
    "/admin",
    "/admin/users",
    "/admin/products",
    "/admin/orders",
    "/admin/categories",
  ]
} satisfies Config;
