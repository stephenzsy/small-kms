import { useMsal } from "@azure/msal-react";
import { Disclosure, Menu, Transition } from "@headlessui/react";
import { Bars3Icon, XMarkIcon } from "@heroicons/react/24/outline";
import classNames from "classnames";
import type { PropsWithChildren } from "react";
import { Fragment, useMemo } from "react";
import { Link, useMatches } from "react-router-dom";
import { RouteIds } from "./route-constants";

export default function Layout(props: PropsWithChildren<{}>) {
  const { accounts, instance } = useMsal();
  const isAuthenticated = accounts.length > 0;
  const authedAccount = accounts[0];
  const matches = useMatches();
  const isCurrentRouteHome = useMemo(() => {
    return matches.some((match) => match.id === RouteIds.home);
  }, [matches]);

  const isCurrentRouteAdmin = useMemo(() => {
    return matches.some((match) => match.id === RouteIds.admin);
  }, [matches]);

  const isAdmin = useMemo(
    () => !!authedAccount.idTokenClaims?.roles?.includes("App.Admin"),
    [authedAccount]
  );

  if (!isAuthenticated) {
    return (
      <main className="grid min-h-full place-items-center px-6 py-24 sm:py-32 lg:px-8">
        <div className="text-center">
          <h1 className="mt-4 text-3xl font-bold tracking-tight text-gray-900 sm:text-5xl">
            Log in to continue
          </h1>

          <div className="mt-10 flex items-center justify-center gap-x-6">
            <button
              type="button"
              className="rounded-md bg-indigo-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
              onClick={() => {
                instance.loginRedirect();
              }}
            >
              Login with Azure Active Directory
            </button>
          </div>
        </div>
      </main>
    );
  }

  return (
    <>
      {/*
        This example requires updating your template:

        ```
        <html class="h-full bg-gray-100">
        <body class="h-full">
        ```
      */}
      <div className="min-h-full">
        <Disclosure as="nav" className="bg-gray-800">
          {({ open }) => (
            <>
              <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                <div className="flex h-16 items-center justify-between">
                  <div className="flex items-center">
                    <div className="flex-shrink-0 font-semibold text-xl text-white">
                      Small KMS
                    </div>
                    <div className="hidden md:block">
                      <div className="ml-10 flex items-baseline space-x-4">
                        <Link
                          to="/"
                          className={classNames(
                            isCurrentRouteHome
                              ? "bg-gray-900 text-white"
                              : "text-gray-300 hover:bg-gray-700 hover:text-white",
                            "rounded-md px-3 py-2 text-sm font-medium"
                          )}
                          aria-current={isCurrentRouteHome ? "page" : undefined}
                        >
                          Home
                        </Link>
                        {isAdmin && (
                          <Link
                            to="/admin"
                            className={classNames(
                              isCurrentRouteAdmin
                                ? "bg-gray-900 text-white"
                                : "text-gray-300 hover:bg-gray-700 hover:text-white",
                              "rounded-md px-3 py-2 text-sm font-medium"
                            )}
                            aria-current={
                              isCurrentRouteAdmin ? "page" : undefined
                            }
                          >
                            Admin
                          </Link>
                        )}
                      </div>
                    </div>
                  </div>
                  <div className="hidden md:block">
                    <div className="ml-4 flex items-center md:ml-6">
                      {/* Profile dropdown */}
                      <Menu as="div" className="relative ml-3">
                        <div>
                          <Menu.Button className="relative flex max-w-xs items-center rounded-full bg-gray-800 text-sm focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-gray-800 text-white">
                            <span className="absolute -inset-1.5" />
                            <span className="sr-only">Open user menu</span>
                            <span>{authedAccount.username}</span>
                          </Menu.Button>
                        </div>
                        <Transition
                          as={Fragment}
                          enter="transition ease-out duration-100"
                          enterFrom="transform opacity-0 scale-95"
                          enterTo="transform opacity-100 scale-100"
                          leave="transition ease-in duration-75"
                          leaveFrom="transform opacity-100 scale-100"
                          leaveTo="transform opacity-0 scale-95"
                        >
                          <Menu.Items className="absolute right-0 z-10 mt-2 w-48 origin-top-right rounded-md bg-white py-1 shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                            <Menu.Item>
                              {({ active }) => (
                                <button
                                  type="button"
                                  onClick={() => {
                                    instance.logoutRedirect();
                                  }}
                                  className={classNames(
                                    active ? "bg-gray-100" : "",
                                    "block px-4 py-2 text-sm text-gray-700"
                                  )}
                                >
                                  Log out
                                </button>
                              )}
                            </Menu.Item>
                          </Menu.Items>
                        </Transition>
                      </Menu>
                    </div>
                  </div>
                  <div className="-mr-2 flex md:hidden">
                    {/* Mobile menu button */}
                    <Disclosure.Button className="relative inline-flex items-center justify-center rounded-md bg-gray-800 p-2 text-gray-400 hover:bg-gray-700 hover:text-white focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-gray-800">
                      <span className="absolute -inset-0.5" />
                      <span className="sr-only">Open main menu</span>
                      {open ? (
                        <XMarkIcon
                          className="block h-6 w-6"
                          aria-hidden="true"
                        />
                      ) : (
                        <Bars3Icon
                          className="block h-6 w-6"
                          aria-hidden="true"
                        />
                      )}
                    </Disclosure.Button>
                  </div>
                </div>
              </div>

              <Disclosure.Panel className="md:hidden">
                <div className="space-y-1 px-2 pb-3 pt-2 sm:px-3">
                  <Disclosure.Button
                    as={Link}
                    to="/"
                    className={classNames(
                      isCurrentRouteHome
                        ? "bg-gray-900 text-white"
                        : "text-gray-300 hover:bg-gray-700 hover:text-white",
                      "block rounded-md px-3 py-2 text-base font-medium"
                    )}
                    aria-current={isCurrentRouteHome ? "page" : undefined}
                  >
                    Home
                  </Disclosure.Button>
                  {isAdmin && (
                    <Disclosure.Button
                      as={Link}
                      to="/admin"
                      className={classNames(
                        isCurrentRouteAdmin
                          ? "bg-gray-900 text-white"
                          : "text-gray-300 hover:bg-gray-700 hover:text-white",
                        "block rounded-md px-3 py-2 text-base font-medium"
                      )}
                      aria-current={isCurrentRouteAdmin ? "page" : undefined}
                    >
                      Admin
                    </Disclosure.Button>
                  )}
                </div>
                <div className="border-t border-gray-700 pb-3 pt-4">
                  <div className="px-5">
                    <div className="text-base font-medium leading-none text-white">
                      {authedAccount.name}
                    </div>
                    <div className="text-sm mt-4 font-medium leading-none text-gray-400">
                      {authedAccount.username}
                    </div>
                  </div>
                  <div className="mt-3 space-y-1 px-2">
                    <Disclosure.Button
                      as="button"
                      onClick={() => {
                        instance.logoutRedirect();
                      }}
                      className="block rounded-md px-3 py-2 text-base font-medium text-gray-400 hover:bg-gray-700 hover:text-white"
                    >
                      Log out
                    </Disclosure.Button>
                  </div>
                </div>
              </Disclosure.Panel>
            </>
          )}
        </Disclosure>

        {props.children}
      </div>
    </>
  );
}
