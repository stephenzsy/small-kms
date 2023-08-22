import { Outlet, RouterProvider, createBrowserRouter } from "react-router-dom";
import AdminLayout from "./admin/layout";
import AdminPage from "./admin/page";
import { LoginProvider } from "./components/auth/LoginProvider";

const router = createBrowserRouter([
  {
    path: "/admin",
    element: (
      <LoginProvider>
        <AdminLayout>
          <Outlet />
        </AdminLayout>
      </LoginProvider>
    ),
    children: [{ index: true, element: <AdminPage /> }],
  },
]);

function App() {
  return <RouterProvider router={router} />;
}

export default App;
