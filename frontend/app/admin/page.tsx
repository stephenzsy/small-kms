import { ChevronRightIcon } from "@heroicons/react/24/outline";
import Link from "next/link";

export default function AdminIndex() {
  return (
    <ul role="list" className="divide-y divide-gray-100">
      <li className="relative gap-x-6 py-5">
        <Link
          className="flex items-center justify-between font-medium"
          href="/admin/ca"
        >
          <span>Manage Certificate Authorities (CA)</span>
          <ChevronRightIcon
            className="h-5 w-5 flex-none text-gray-400"
            aria-hidden="true"
          />
        </Link>
      </li>
    </ul>
  );
}
