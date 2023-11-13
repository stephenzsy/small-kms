export const dateShortFormatter = new Intl.DateTimeFormat("en-US", {
  year: "numeric",
  month: "2-digit",
  day: "2-digit",
});

export function ShortDate({ numericDate }: { numericDate?: number }) {
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
