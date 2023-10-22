import React from "react";
import { createBrowserRouter, Outlet } from "react-router-dom";
import { AdminLayout } from "./admin/Layout";
import { NamespaceContextProvider } from "./admin/NamespaceContext";
import AdminPage from "./admin/Page";
import { AuthProvider } from "./auth/AuthProvider";
import AppLayout from "./Layout";
import { RouteIds } from "./route-constants";

const DiagnosticsPage = React.lazy(() => import("./diagnostics/Page"));

const MainPage = React.lazy(() => import("./MainPage"));
const AdminEnrollPage = React.lazy(() => import("./admin/AdminEnroll"));
const NamespacePage = React.lazy(() => import("./admin/NamespacePage"));
const CertificateTemplatePage = React.lazy(
  () => import("./admin/CertificateTemplatePage")
);
const CertificatePage = React.lazy(() => import("./admin/CertificatePage"));
const ServicePage = React.lazy(() => import("./service/Page"));
const AgentDashboardPage = React.lazy(
  () => import("./admin/AgentDashboardPage")
);
const ManagedAppsPage = React.lazy(() => import("./admin/ManagedAppsPage"));
const RegisterPage = React.lazy(() => import("./admin/RegisterPage"));
const CAsPage = React.lazy(() => import("./admin/CAsPage"));

export const router = createBrowserRouter([
  {
    path: "/",
    element: (
      <AuthProvider>
        <AppLayout>
          <React.Suspense>
            <Outlet />
          </React.Suspense>
        </AppLayout>
      </AuthProvider>
    ),
    children: [
      { index: true, element: <MainPage />, id: RouteIds.home },
      { path: "diagnostics", element: <DiagnosticsPage /> },
      {
        element: (
          <AdminLayout>
            <Outlet />
          </AdminLayout>
        ),
        children: [
          {
            path: "apps",
            id: RouteIds.apps,
            element: <Outlet />,
            children: [{ index: true, element: <ManagedAppsPage /> }],
          },
          {
            path: "cas",
            element: <Outlet />,
            children: [{ index: true, element: <CAsPage /> }],
          },
          {
            path: "admin",
            id: RouteIds.admin,
            element: <Outlet />,
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
                        ],
                      },
                      {
                        path: "certificates/:certId",
                        element: <CertificatePage />,
                      },
                      {
                        path: "agent",
                        element: <AgentDashboardPage />,
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
                path: "settings",
                element: <ServicePage />,
              },
              {
                path: "enroll",
                element: <AdminEnrollPage />,
              },
            ],
          },
        ],
      },
    ],
  },
]);
