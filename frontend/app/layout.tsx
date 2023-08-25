import { getMsAuth } from "@/utils/aadAuthUtils";
import classNames from "classnames";
import type { Metadata } from "next";
import { Inter } from "next/font/google";
import LayoutDisclosure from "./_layoutDisclosure";
import "./globals.css";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Small KMS",
  description: "Key management system",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const auth = getMsAuth();
  return (
    <html lang="en" className="h-full">
      <body className={classNames(inter.className, "h-full")}>
        <div className="min-h-full">
          <LayoutDisclosure authClient={auth.client} />
          <div className="py-10">{children}</div>
        </div>
      </body>
    </html>
  );
}
