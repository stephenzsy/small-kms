import React from "react";
import { createBrowserRouter, Outlet } from "react-router-dom";
import { AdminLayout } from "./admin/Layout";
import { NamespaceContextProvider } from "./admin/NamespaceContext";
import AdminPage from "./admin/Page";
import RegisterPage from "./admin/RegisterPage";
import { AuthProvider } from "./auth/AuthProvider";
import Layout from "./Layout";
import { RouteIds } from "./route-constants";

const DiagnosticsPage = React.lazy(() => import("./diagnostics/Page"));

const MainPage = React.lazy(() => import("./MainPage"));
const AdminEnrollPage = React.lazy(() => import("./admin/AdminEnroll"));
const NamespacePage = React.lazy(() => import("./admin/NamespacePage"));
const CertificateTemplatePage = React.lazy(
  () => import("./admin/CertificateTemplatePage")
);

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
            path: ":profileType",
            children: [
              {
                path: ":namespaceId",
                element: (
                  <NamespaceContextProvider>
                    <Outlet />
                  </NamespaceContextProvider>
                ),
                children: [
                  { index: true, element: <NamespacePage /> },
                  {
                    path: "certificate-templates/:templateId",
                    children: [
                      { index: true, element: <CertificateTemplatePage /> },
                      {
                        path: "certificates/:certId",
                      },
                    ],
                  },
                ],
              },
              {
                path: "register",
                element: <RegisterPage />,
              },
            ],
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
