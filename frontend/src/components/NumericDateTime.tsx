
export function NumericDateTime({
    value, format,
}: {
    value?: number;
    format?: (date: Date) => string;
}) {
    if (!value) {
        return null;
    }
    const ts = new Date(value * 1000);
    return (
        <time dateTime={ts.toISOString()} className="tabular-nums">
            {format ? format(ts) : ts.toLocaleString()}
        </time>
    );
}
