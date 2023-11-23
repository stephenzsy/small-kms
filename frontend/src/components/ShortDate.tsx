import { dateShortFormatter } from "../utils/datetimeUtils";


export function ShortDate({ numericDate }: { numericDate?: number; }) {
  if (numericDate) {
    const ts = new Date(numericDate * 1000);
    return (
      <time className="tabular-nums" dateTime={ts.toISOString()}>
        {dateShortFormatter.format(ts)}
      </time>
    );
  }
  return null;
}
