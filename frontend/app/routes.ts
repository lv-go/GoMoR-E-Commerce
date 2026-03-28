import { type RouteConfig, index, layout, prefix, route } from "@react-router/dev/routes";

export default [
  index("routes/home.tsx"),
  route("shop", "routes/shop.tsx"),
  // route("product/:id", "routes/product/:id.tsx"),
  route("cart", "routes/cart.tsx"),
  route("favorites", "routes/favorites.tsx"),
  route("profile", "routes/user/profile.tsx"),
  route("login", "routes/auth/login.tsx"),
  route("register", "routes/auth/register.tsx"),
  layout("routes/admin/layout.tsx", prefix("admin", [
    index("routes/admin/admin-dashboard.tsx"),
    route("users", "routes/admin/user-list.tsx"),
    route("products", "routes/admin/product-list.tsx"),
    route("orders", "routes/admin/order-list.tsx"),
    route("categories", "routes/admin/category-list.tsx"),
  ])),
] satisfies RouteConfig;
