import React from "react";
import ReactDOM from "react-dom/client";
import { RouterProvider } from "react-router-dom";
import "./index.css";
import { router } from "./routes";
import { App } from "antd";
ReactDOM.createRoot(document.getElementById("root")!).render(
  <App>
    <RouterProvider router={router} />
  </App>
);
