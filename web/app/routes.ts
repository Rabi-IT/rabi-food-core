import { type RouteConfig, index, layout, route } from "@react-router/dev/routes";

export default [
  layout("routes/storefront/layout.tsx", [
    index("routes/storefront/storefront-page.tsx"),
    route("product/:id", "routes/storefront/product.tsx"),
  ]),
] satisfies RouteConfig;
