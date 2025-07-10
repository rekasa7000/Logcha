import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { useAuth } from "@/lib/auth-context";
import { Button } from "@/components/ui/button";
import { ThemeToggleButton } from "@/components/shared/theme-button";
import { LogOut } from "lucide-react";

export const Route = createFileRoute("/_protected")({
  component: RouteComponent,
  beforeLoad: () => {
    const token = localStorage.getItem('token');
    if (!token) {
      throw redirect({
        to: '/signin',
      });
    }
  },
});

function RouteComponent() {
  const { user, logout, isLoading } = useAuth();

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
          <p className="text-muted-foreground">Loading...</p>
        </div>
      </div>
    );
  }

  if (!user) {
    return null;
  }

  return (
    <div className="min-h-screen bg-background">
      <header className="border-b">
        <div className="container mx-auto px-4 py-3 flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <img
              src="/logo-transparent.png"
              alt="Logcha"
              className="h-8 w-auto dark:invert"
            />
            <div className="text-sm text-muted-foreground">
              Welcome, {user.first_name} {user.last_name}
            </div>
          </div>
          <div className="flex items-center space-x-2">
            <ThemeToggleButton />
            <Button
              variant="outline"
              size="sm"
              onClick={logout}
              className="flex items-center space-x-2"
            >
              <LogOut className="h-4 w-4" />
              <span>Logout</span>
            </Button>
          </div>
        </div>
      </header>
      <main className="container mx-auto px-4 py-6">
        <Outlet />
      </main>
    </div>
  );
}