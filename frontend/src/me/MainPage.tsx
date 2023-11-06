import { Card, Typography } from "antd";

export default function MainPage(props: React.PropsWithChildren<{}>) {
  return (
    <main className="p-6 max-w-7xl mx-auto space-y-6">{props.children}</main>
  );
}
