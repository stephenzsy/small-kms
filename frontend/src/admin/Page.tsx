import { ChevronRightIcon } from "@heroicons/react/24/outline";
import { Link } from "react-router-dom";

export default function AdminPage() {
  return (
    <>
      <section className="divide-y divide-gray-200 overflow-hidden rounded-lg bg-white shadow">
        <h2 className="px-4 py-5 sm:px-6 text-lg font-semibold">Actions</h2>
        <div className="px-4 py-5 sm:p-6">
          <Link to="./ca" className="inline-flex items-center gap-x-ex">
            <span>Manage Certificate Authorities</span>
            <ChevronRightIcon className="h-em w-em" />
          </Link>
        </div>
        <div className="px-4 py-5 sm:p-6">
          <Link to="./testca" className="inline-flex items-center gap-x-ex">
            <span>Manage Certificate Authorities (Test)</span>
            <ChevronRightIcon className="h-em w-em" />
          </Link>
        </div>
      </section>
    </>
  );
}
