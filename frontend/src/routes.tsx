import { createBrowserRouter, Outlet } from "react-router-dom";
import { AdminLayout } from "./admin/Layout";
import AdminPage from "./admin/Page";
import PoliciesPage from "./admin/PoliciesPage";
import { AuthProvider } from "./auth/AuthProvider";
import Layout from "./Layout";
import { RouteIds } from "./route-constants";
import RegisterPage from "./admin/RegisterPage";
import React from "react";

const DiagnosticsPage = React.lazy(() => import("./diagnostics/Page"));
const MainPage = React.lazy(() => import("./MainPage"));
const AdminEnrollPage = React.lazy(() => import("./admin/AdminEnroll"));
const CertificatePage = React.lazy(() => import("./admin/CertificatePage"));
const PolicyPage = React.lazy(() => import("./admin/PolicyPage"));
const PermissionsPage = React.lazy(() => import("./admin/PermissionsPage"));

export const router = createBrowserRouter([
  {
    path: "/",
    element: (
      <AuthProvider>
        <Layout>
          <React.Suspense>
            <Outlet />
          </React.Suspense>
        </Layout>
      </AuthProvider>
    ),
    children: [
      { index: true, element: <MainPage />, id: RouteIds.home },
      { path: "diagnostics", element: <DiagnosticsPage /> },
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
              {
                path: ":policyId/latest-certificate",
                element: <CertificatePage />,
              },
            ],
          },
          {
            path: ":namespaceId/permissions",
            children: [
              { index: true, element: <PermissionsPage /> },
              { path: ":policyId", element: <PolicyPage /> },
              {
                path: ":policyId/latest-certificate",
                element: <CertificatePage />,
              },
            ],
          },
          {
            path: "register",
            element: <RegisterPage />,
          },
          {
            path: "enroll",
            element: <AdminEnrollPage />,
          },
        ],
      },
    ],
  },
]);
