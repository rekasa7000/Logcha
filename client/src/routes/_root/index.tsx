import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_root/")({
  component: App,
});

function App() {
  return <div className="text-center">Hello World</div>;
}
