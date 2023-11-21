import React from "react";
import { createBrowserRouter, Outlet } from "react-router-dom";
import { ManagedAppContextProvider } from "./admin/contexts/ManagedAppContext";
import { NamespaceContext } from "./admin/contexts/NamespaceContext";
import {
  NamespaceContextRouteProvider,
  NamespaceContextRouteProvider2,
} from "./admin/contexts/NamespaceContextRouteProvider";
import { AdminLayout } from "./admin/Layout";
import { AuthProvider } from "./auth/AuthProvider";
import { NamespaceKind } from "./generated";
import { NamespaceProvider } from "./generated/apiv2";
import AppLayout from "./Layout";
import { RouteIds } from "./route-constants";

const ProfilesPage = React.lazy(() => import("./admin/ProfilesPage"));
const AgentPage = React.lazy(() => import("./agents/[id]/page"));
const AgentsPage = React.lazy(() => import("./agents/page"));
const SystemAppsPage = React.lazy(() => import("./system/page"));
const KeyPage = React.lazy(() => import("./admin/KeyPage"));
const KeyPolicyPage = React.lazy(() => import("./admin/KeyPolicyPage"));
const KeyPolicyPage2 = React.lazy(() => import("./key-policies/[id]/page"));
const SecretPage = React.lazy(() => import("./admin/SecretPage"));
const SecretPolicyPage = React.lazy(() => import("./admin/SecretPolicyPage"));
const DiagnosticsPage = React.lazy(() => import("./diagnostics/Page"));
const MainPage = React.lazy(() => import("./me/MainPage"));
const NamespacePage = React.lazy(() => import("./admin/NamespacePage"));
const NamespacePage2 = React.lazy(() => import("./admin/NamespacePage2"));
const CertificatePage = React.lazy(() => import("./certificates/[id]/page"));
const AppsPage = React.lazy(() => import("./admin/AppsPage"));
const CAsPage = React.lazy(() => import("./admin/CAsPage"));
const CertPolicyPage = React.lazy(() => import("./admin/CertPolicyPage"));
const CertPolicyPage2 = React.lazy(() => import("./cert-policies/[id]/page"));
const ManagedAppPage = React.lazy(() => import("./admin/ManagedAppPage"));
const ProvisionAgentPage = React.lazy(
  () => import("./admin/ProvisionAgentPage")
);
const AgentDashboardPage = React.lazy(
  () => import("./admin/AgentDashboardPage")
);
const RadiusConfigPage = React.lazy(() => import("./admin/RadiusConfigPage"));

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
      {
        path: "agents",
        children: [
          {
            index: true,
            element: <AgentsPage />,
          },
          {
            path: ":id",
            element: <AgentPage />,
          },
        ],
      },
      { path: "system", element: <SystemAppsPage /> },
      {
        path: ":nsKind/:nsId",
        element: (
          <NamespaceContextRouteProvider2>
            <Outlet />
          </NamespaceContextRouteProvider2>
        ),
        children: [
          { index: true, element: <NamespacePage2 /> },
          { path: "cert-policies/:id", element: <CertPolicyPage2 /> },
          { path: "certificates/:id", element: <CertificatePage /> },
          { path: "key-policies/:id", element: <KeyPolicyPage2 /> },
        ],
      },
      {
        path: "groups",
        element: (
          <ProfilesPage
            title="Groups"
            namespaceProvider={NamespaceProvider.NamespaceProviderGroup}
          />
        ),
      },
      {
        path: "users",
        element: (
          <ProfilesPage
            title="Users"
            namespaceProvider={NamespaceProvider.NamespaceProviderUser}
          />
        ),
      },
      {
        path: "service-principals",
        element: (
          <ProfilesPage
            title="Service Principals"
            namespaceProvider={
              NamespaceProvider.NamespaceProviderServicePrincipal
            }
          />
        ),
      },
      {
        path: "entra/:nsKind",
        children: [
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
                  <ManagedAppContextProvider>
                    <ManagedAppPage isSystemApp />
                  </ManagedAppContextProvider>
                ),
              },
              {
                path: "system/default/provision-agent",
                element: <ProvisionAgentPage isGlobalConfig />,
              },
              {
                path: "system/default/radius-config",
                element: <RadiusConfigPage isGlobalConfig />,
              },
              {
                path: "managed/:appId",
                element: (
                  <ManagedAppContextProvider>
                    <Outlet />
                  </ManagedAppContextProvider>
                ),
                children: [
                  { index: true, element: <ManagedAppPage /> },
                  { path: "provision-agent", element: <ProvisionAgentPage /> },
                  { path: "radius-config", element: <RadiusConfigPage /> },
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
                    path: "key-policies/:policyId",
                    element: <KeyPolicyPage />,
                  },
                  {
                    path: "secret-policy/:policyId",
                    element: <SecretPolicyPage />,
                  },
                  {
                    path: "secrets/:id",
                    element: <SecretPage />,
                  },
                  {
                    path: "cert/:certId",
                    element: <CertificatePage />,
                  },
                  {
                    path: "keys/:id",
                    element: <KeyPage />,
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
        ],
      },
    ],
  },
]);
