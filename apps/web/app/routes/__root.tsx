import { createRootRoute, Outlet } from '@tanstack/react-router';
import { TanStackRouterDevtools } from '@tanstack/router-devtools';

export const Route = createRootRoute({
  component: RootComponent,
});

function RootComponent() {
  return (
    <div className="bg-background min-h-screen font-sans antialiased">
      <Outlet />
      {import.meta.env.DEV && <TanStackRouterDevtools />}
    </div>
  );
}
