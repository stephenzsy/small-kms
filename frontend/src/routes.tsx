import { createBrowserRouter, Outlet } from "react-router-dom";
import { AdminLayout } from "./admin/Layout";
import AdminPage from "./admin/Page";
import PoliciesPage from "./admin/PoliciesPage";
import PolicyPage from "./admin/PolicyPage";
import { AuthProvider } from "./auth/AuthProvider";
import Layout from "./Layout";
import { MainPage } from "./MainPage";
import { RouteIds } from "./route-constants";
import RegisterPage from "./admin/RegisterPage";

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
          {
            path: ":namespaceId/policies",
            children: [
              { index: true, element: <PoliciesPage /> },
              { path: ":policyId", element: <PolicyPage /> },
            ],
          },
          {
            path: "register",
            element: <RegisterPage />,
          },
        ],
      },
    ],
  },
]);
