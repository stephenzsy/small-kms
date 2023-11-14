"use client";

import { useMsal } from "@azure/msal-react";
import { useEffect } from "react";

export default function Logout() {
  const { instance } = useMsal();
  useEffect(() => {
    instance.logoutRedirect();
  });
  return (
    <div>
      <h1>Logging you out</h1>
    </div>
  );
}
