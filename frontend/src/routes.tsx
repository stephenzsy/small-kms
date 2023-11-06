import React from "react";
import { createBrowserRouter, Outlet } from "react-router-dom";
import { ManagedAppContextProvider } from "./admin/contexts/ManagedAppContext";
import { NamespaceContextRouteProvider } from "./admin/contexts/NamespaceContextRouteProvider";
import { AdminLayout } from "./admin/Layout";
import { AuthProvider } from "./auth/AuthProvider";
import AppLayout from "./Layout";
import { RouteIds } from "./route-constants";
import { NamespaceContext } from "./admin/contexts/NamespaceContext";
import { NamespaceKind } from "./generated";

const EntraProfilesPage = React.lazy(() => import("./admin/EntraProfilesPage"));
const SecretPolicyPage = React.lazy(() => import("./admin/SecretPolicyPage"));
const DiagnosticsPage = React.lazy(() => import("./diagnostics/Page"));
const MainPage = React.lazy(() => import("./me/MainPage"));
const NamespacePage = React.lazy(() => import("./admin/NamespacePage"));
const CertificatePage = React.lazy(() => import("./admin/CertificatePage"));
const AppsPage = React.lazy(() => import("./admin/AppsPage"));
const CAsPage = React.lazy(() => import("./admin/CAsPage"));
const CertPolicyPage = React.lazy(() => import("./admin/CertPolicyPage"));
const ManagedAppPage = React.lazy(() => import("./admin/ManagedAppPage"));
const ProvisionAgentPage = React.lazy(
  () => import("./admin/ProvisionAgentPage")
);

const AgentDashboardPage = React.lazy(
  () => import("./admin/AgentDashboardPage")
);

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
      {
        index: true,
        element: (
          <NamespaceContext.Provider
            value={{
              namespaceId: "me",
              namespaceKind: NamespaceKind.NamespaceKindUser,
            }}
          >
            <MainPage>
              <NamespacePage />
            </MainPage>
          </NamespaceContext.Provider>
        ),
        id: RouteIds.home,
      },
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
            children: [
              { index: true, element: <AppsPage /> },
              {
                path: "system/:appId",
                element: (
                  <ManagedAppContextProvider isSystemApp={true}>
                    <ManagedAppPage isSystemApp />
                  </ManagedAppContextProvider>
                ),
              },
              {
                path: "system/default/provision-agent",
                element: <ProvisionAgentPage isGlobalConfig />,
              },
              {
                path: "managed/:appId",
                element: (
                  <ManagedAppContextProvider isSystemApp={false}>
                    <Outlet />
                  </ManagedAppContextProvider>
                ),
                children: [
                  { index: true, element: <ManagedAppPage /> },
                  { path: "provision-agent", element: <ProvisionAgentPage /> },
                ],
              },
              {
                path: ":nsKind/:nsId",
                element: (
                  <NamespaceContextRouteProvider>
                    <Outlet />
                  </NamespaceContextRouteProvider>
                ),
                children: [
                  { index: true, element: <NamespacePage /> },
                  {
                    path: "cert-policy/:certPolicyId",
                    element: <CertPolicyPage />,
                  },
                  {
                    path: "secret-policy/:policyId",
                    element: <SecretPolicyPage />,
                  },
                  {
                    path: "cert/:certId",
                    element: <CertificatePage />,
                  },
                  {
                    path: "agent/:instanceId/dashboard",
                    element: <AgentDashboardPage />,
                  },
                ],
              },
            ],
          },
          {
            path: "ca",
            element: <Outlet />,
            children: [
              { index: true, element: <CAsPage /> },
              {
                path: ":nsKind/:nsId",
                element: (
                  <NamespaceContextRouteProvider>
                    <Outlet />
                  </NamespaceContextRouteProvider>
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
            path: "entra/:nsKind",
            children: [
              {
                index: true,
                element: <EntraProfilesPage />,
              },
              {
                path: ":nsId",
                element: (
                  <NamespaceContextRouteProvider>
                    <Outlet />
                  </NamespaceContextRouteProvider>
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
        ],
      },
    ],
  },
]);
