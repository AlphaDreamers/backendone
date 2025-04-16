import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import SharedLayout from "./components/shared-layout";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "BlockChain Service",
  description: "A platform for blockchain services and transactions",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <SharedLayout>
          {children}
        </SharedLayout>
      </body>
    </html>
  );
}
