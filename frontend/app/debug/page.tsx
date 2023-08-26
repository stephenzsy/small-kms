import _ from "lodash";
import { headers } from "next/headers";

export default function DebugPage() {
  const allHeaders = _.toArray(headers().entries);
  return <pre>{JSON.stringify(allHeaders, null, 2)}</pre>;
}
