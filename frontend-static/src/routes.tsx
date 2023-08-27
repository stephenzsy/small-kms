import { createBrowserRouter, Outlet } from "react-router-dom";
import Layout from "./Layout";
import { MainPage } from "./MainPage";
import { AuthProvider } from "./auth/AuthProvider";
import { RouteIds } from "./route-constants";
import { AdminLayout } from "./admin/Layout";
import AdminPage from "./admin/Page";
import AdminCaPage from "./admin/CaPage";
import CreateCertPage from "./admin/CreateCertPage";

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
            <Outlet />
          </AdminLayout>
        ),
        children: [
          { index: true, element: <AdminPage /> },
          { path: "ca", element: <AdminCaPage /> },
          { path: "cert/:namespaceId/new", element: <CreateCertPage /> },
        ],
      },
    ],
  },
]);