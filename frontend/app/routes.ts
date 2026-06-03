import { index, type RouteConfig, route } from "@react-router/dev/routes";

export default [
  index("routes/home.tsx"),
  route("login", "routes/login.tsx"),
  route("registration", "routes/registration.tsx"),
  route("recovery", "routes/recovery.tsx"),
  route("verification", "routes/verification.tsx"),
  route("dashboard", "routes/dashboard.tsx"),
  route("settings", "routes/settings.tsx"),
  route("users/:username/profile", "routes/users.$username.profile.tsx"),
  route("sources/books/:id", "routes/sources.books.$id.tsx"),
] satisfies RouteConfig;
