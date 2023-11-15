import {
  Disclosure,
  Menu as HeadlessUIMenu,
  Transition,
} from "@headlessui/react";
import {
  ArrowRightOnRectangleIcon,
  Bars3Icon,
  UserCircleIcon,
  UserIcon,
  XMarkIcon,
} from "@heroicons/react/24/outline";
import { Avatar, Button, Layout, Menu, MenuProps } from "antd";
import classNames from "classnames";
import type { PropsWithChildren } from "react";
import { Fragment, useMemo } from "react";
import { Link, NavLink, useMatches } from "react-router-dom";
import { useAppAuthContext } from "./auth/AuthProvider";
import { RouteIds } from "./route-constants";

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
  const { account, isAuthenticated, logout } = useAppAuthContext();
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
export default function AppLayout(props: PropsWithChildren<{}>) {
  const { account, logout, isAuthenticated } = useAppAuthContext();
  const matches = useMatches();
  const isCurrentRouteHome = useMemo(() => {
    return matches.some((match) => match.id === RouteIds.home);
  }, [matches]);

  const isCurrentRouteAdmin = useMemo(() => {
    return matches.some((match) => match.id === RouteIds.admin);
  }, [matches]);

  const isAdmin = useMemo(
    () => !!account?.idTokenClaims?.roles?.includes("App.Admin"),
    [account]
  );

  const navItems = useNavItems(isAdmin);
  const userNavItems = useUserMenuItems();

  return (
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
        {props.children}
      </Layout.Content>
    </Layout>
  );
}
