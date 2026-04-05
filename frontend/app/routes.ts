import { type RouteConfig, index, layout, prefix, route } from "@react-router/dev/routes";

export default [
  index("routes/home.tsx"),
  route("shop", "routes/shop.tsx"),
  route("product/:id", "routes/product-details.tsx"),
  route("cart", "routes/cart.tsx"),
  route("favorites", "routes/favorites.tsx"),
  route("profile", "routes/user/profile.tsx"),
  route("user-orders", "routes/user/orders.tsx"),
  route("login", "routes/auth/login.tsx"),
  route("register", "routes/auth/register.tsx"),
  layout("routes/admin/layout.tsx", prefix("admin", [
    index("routes/admin/dashboard.tsx"),
    route("users", "routes/admin/users.tsx"),
    route("products", "routes/admin/products.tsx"),
    route("orders", "routes/admin/orders.tsx"),
    route("categories", "routes/admin/categories.tsx"),
  ])),
] satisfies RouteConfig;
