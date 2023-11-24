import { useBoolean, useMemoizedFn } from "ahooks";
import { Drawer, DrawerProps } from "antd";
import React, { PropsWithChildren } from "react";

type DrawerContextValue = {
  openDrawer: (
    content: React.ReactNode,
    props?: Pick<DrawerProps, "title" | "size">
  ) => void;
};

export const DrawerContext = React.createContext<DrawerContextValue>({
  openDrawer: () => {},
});

export function DrawerProvider({ children }: PropsWithChildren) {
  const [drawerOpen, { setTrue: setDrawerOpenTrue, setFalse: closeDrawer }] =
    useBoolean(false);
  const [{ children: drawerChildren, ...restProps }, setDrawerProps] =
    React.useState<Pick<DrawerProps, "title" | "children">>({});

  const openDrawer = useMemoizedFn(
    (content: React.ReactNode, props?: Pick<DrawerProps, "title" | "size">) => {
      setDrawerOpenTrue();
      setDrawerProps({
        ...props,
        children: content,
      });
    }
  );

  return (
    <DrawerContext.Provider
      value={{
        openDrawer,
      }}
    >
      {children}
      <Drawer open={drawerOpen} onClose={closeDrawer} {...restProps}>
        {drawerOpen && drawerChildren}
      </Drawer>
    </DrawerContext.Provider>
  );
}
