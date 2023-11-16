import {
  AuthenticatedTemplate,
  UnauthenticatedTemplate,
} from "@azure/msal-react";
import {
  ArrowLeftOnRectangleIcon,
  ArrowRightOnRectangleIcon,
  UserCircleIcon,
} from "@heroicons/react/24/outline";
import { Avatar, Button, ConfigProvider, Layout, Menu, MenuProps } from "antd";
import type { PropsWithChildren } from "react";
import { useMemo } from "react";
import { NavLink, useMatches } from "react-router-dom";
import { useAppAuthContext } from "./auth/AuthProvider";

function useNavItems(isAdmin: boolean): MenuProps["items"] {
  return useMemo(
    () =>
      isAdmin
        ? [
            {
              key: "agents",
              label: <NavLink to="/agents">Agents</NavLink>,
            },
            {
              key: "system",
              label: <NavLink to="/system">System</NavLink>,
            },
            {
              key: "ca",
              label: <NavLink to="/ca">CAs</NavLink>,
            },
            {
              key: "app",
              label: <NavLink to="/app">Apps</NavLink>,
            },
            {
              key: "groups",
              label: <NavLink to="/entra/group">Groups</NavLink>,
            },
            {
              key: "users",
              label: <NavLink to="/entra/user">Users</NavLink>,
            },
          ]
        : [],
    [isAdmin]
  );
}

function useUserMenuItems(): MenuProps["items"] {
  const { account, isAuthenticated, logout, login } = useAppAuthContext();
  console.log;
  return useMemo(
    () =>
      isAuthenticated
        ? [
            {
              key: "user",
              label: (
                <Avatar icon={<UserCircleIcon className="h-full w-full" />} />
              ),
              children: [
                {
                  key: "authed-user-info",
                  type: "group",
                  label: (
                    <div className="cursor-default">
                      <div className="text-white">{account?.name}</div>
                      <div>{account?.username}</div>
                    </div>
                  ),
                  children: [
                    {
                      key: "logout",
                      icon: <ArrowRightOnRectangleIcon className="h-4 w-4" />,
                      label: "Logout",
                      onClick: logout,
                    },
                  ],
                },
              ],
            },
          ]
        : [
            {
              key: "login",
              label: <Button>Login</Button>,
            },
          ],
    [isAuthenticated, account]
  );
}

const theme = {
  token: {
    fontFamily: "Mona Sans, ui-sans-serif, system-ui, sans-serif",
  },
};

export default function AppLayout(props: PropsWithChildren<{}>) {
  const { account, login } = useAppAuthContext();
  const matches = useMatches();
  const isAdmin = useMemo(
    () => !!account?.idTokenClaims?.roles?.includes("App.Admin"),
    [account]
  );

  const navItems = useNavItems(isAdmin);
  const userNavItems = useUserMenuItems();

  return (
    <ConfigProvider theme={theme}>
      <Layout className="min-h-full">
        <Layout.Header className="flex items-center gap-6">
          <NavLink to="/" className="text-2xl text-white">
            CryptoCat
          </NavLink>
          <Menu
            className="flex-auto"
            theme="dark"
            mode="horizontal"
            items={navItems}
          />
          <Menu items={userNavItems} theme="dark" mode="horizontal" />
        </Layout.Header>
        <Layout.Content className="p-6 max-w-7xl mx-auto w-full space-y-6">
          <AuthenticatedTemplate>{props.children}</AuthenticatedTemplate>
          <UnauthenticatedTemplate>
            <center className="mt-10 lg:mt-40 space-y-4">
              <div className="text-4xl">
                You must be logged in to access this app
              </div>
              <Button
                type="primary"
                className="inline-flex items-center"
                icon={<ArrowLeftOnRectangleIcon className="h-4 w-4" />}
                onClick={login}
              >
                Login
              </Button>
            </center>
          </UnauthenticatedTemplate>
        </Layout.Content>
      </Layout>
    </ConfigProvider>
  );
}
