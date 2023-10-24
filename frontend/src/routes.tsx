import React from "react";
import { createBrowserRouter, Outlet } from "react-router-dom";
import { AdminLayout } from "./admin/Layout";
import { NamespaceContextProvider } from "./admin/NamespaceContext";
import { NamespaceContextProvider as NamespaceContextProvider2 } from "./admin/NamespaceContext2";
import AdminPage from "./admin/Page";
import { AuthProvider } from "./auth/AuthProvider";
import AppLayout from "./Layout";
import { RouteIds } from "./route-constants";

const DiagnosticsPage = React.lazy(() => import("./diagnostics/Page"));

const MainPage = React.lazy(() => import("./MainPage"));
const AdminEnrollPage = React.lazy(() => import("./admin/AdminEnroll"));
const NamespacePage = React.lazy(() => import("./admin/NamespacePage"));
const CertificatePage = React.lazy(() => import("./admin/CertificatePage"));
const ServicePage = React.lazy(() => import("./service/Page"));
const AgentDashboardPage = React.lazy(
  () => import("./admin/AgentDashboardPage")
);
const AppsPage = React.lazy(() => import("./admin/AppsPage"));
const CAsPage = React.lazy(() => import("./admin/CAsPage"));
const CertPolicyPage = React.lazy(() => import("./admin/CertPolicyPage"));

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
            path: "app",
            id: RouteIds.apps,
            element: <Outlet />,
            children: [{ index: true, element: <AppsPage /> }],
          },
          {
            path: "ca",
            element: <Outlet />,
            children: [
              { index: true, element: <CAsPage /> },
              {
                path: ":nsKind/:nsId",
                element: (
                  <NamespaceContextProvider2>
                    <Outlet />
                  </NamespaceContextProvider2>
                ),
                children: [
                  { index: true, element: <NamespacePage /> },
                  {
                    path: "cert-policy/:certPolicyId",
                    element: <CertPolicyPage />,
                  },
                  {
                    path: "cert/:certId",
                    element: <CertificatePage />,
                  },
                ],
              },
            ],
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
                      {
                        path: "agent",
                        element: <AgentDashboardPage />,
                      },
                    ],
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
