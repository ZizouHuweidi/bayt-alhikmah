import type { Route } from "./+types/home";
import Landing from "./landing";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Bayt Alhikmah - Knowledge Management" },
    { name: "description", content: "Organize and access your knowledge efficiently with Bayt Alhikmah." },
  ];
}

export default function Home() {
  return <Landing />;
}
