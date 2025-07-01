import { ThemeToggleButton } from "@/components/shared/theme-button";
import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth")({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="min-h-screen w-full">
      <ThemeToggleButton />
      <Outlet />
    </div>
  );
}
