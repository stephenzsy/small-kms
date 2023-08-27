import { createBrowserRouter, Outlet } from "react-router-dom";
import Layout from "./Layout";
import { MainPage } from "./MainPage";
import { AuthProvider } from "./auth/AuthProvider";
import { RouteIds } from "./route-constants";
import { AdminLayout } from "./admin/Layout";

export const router = createBrowserRouter([
  {
    path: "/",
    element: (
      <AuthProvider>
        <Layout>
          <Outlet />
        </Layout>
      </AuthProvider>
    ),
    children: [
      { index: true, element: <MainPage />, id: RouteIds.home },
      {
        path: "admin",
        id: RouteIds.admin,
        element: (
          <AdminLayout>
            <div>Admin</div>
          </AdminLayout>
        ),
      },
    ],
  },
]);
