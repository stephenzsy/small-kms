import { Link, type To } from "react-router-dom";

export interface BreadcrumbPageMetadata {
  to: To;
  name: string;
}
export function AdminBreadcrumb({
  pages: pages,
}: {
  pages: BreadcrumbPageMetadata[];
}) {
  return (
    <nav className="flex" aria-label="Breadcrumb">
      <ol
        role="list"
        className="flex w-full space-x-4 rounded-md bg-white px-6 shadow "
      >
        <li className="flex">
          <div className="flex items-center">
            <Link to="/admin" className="text-gray-400 hover:text-gray-500">
              Admin home
            </Link>
          </div>
        </li>
        {pages.map((page, index) => (
          <li key={index} className="flex">
            <div className="flex items-center">
              <svg
                className="h-full w-6 flex-shrink-0 text-gray-200"
                viewBox="0 0 24 44"
                preserveAspectRatio="none"
                fill="currentColor"
                aria-hidden="true"
              >
                <path d="M.293 0l22 22-22 22h1.414l22-22-22-22H.293z" />
              </svg>
              <Link
                to={page.to}
                className="ml-4 text-sm font-medium text-gray-500 hover:text-gray-700"
                aria-current={index + 1 >= pages.length ? "page" : undefined}
              >
                {page.name}
              </Link>
            </div>
          </li>
        ))}
      </ol>
    </nav>
  );
}
