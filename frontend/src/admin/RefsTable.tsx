export type RefTableItem = {
  id: string;
};

export interface RefTableColumn<T extends RefTableItem> {
  columnKey: string;
  header: string;
  render: (item: T) => React.ReactNode;
}

export function RefsTable<T extends RefTableItem>(props: {
  items: T[] | undefined;
  columns?: RefTableColumn<T>[];
  title: string;
  tableActions?: React.ReactNode;
  refActions?: (ref: T) => React.ReactNode;
}) {
  const { title, tableActions, items, refActions, columns = [] } = props;

  return (
    <section className="overflow-hidden rounded-lg bg-white shadow px-4 sm:px-6 lg:px-8 py-6">
      <div className="flex flex-row items-center justify-between">
        <h2 className="text-lg font-semibold">{title}</h2>
        {tableActions}
      </div>
      <div className="mt-8 flow-root">
        <div className="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
          <div className="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
            {items ? (
              items.length === 0 ? (
                <div>No items</div>
              ) : (
                <table className="min-w-full divide-y divide-gray-300">
                  <thead>
                    <tr>
                      <th
                        scope="col"
                        className="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 sm:pl-0"
                      >
                        ID
                      </th>
                      {columns.map((c) => (
                        <th
                          key={c.columnKey}
                          scope="col"
                          className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900"
                        >
                          {c.header}
                        </th>
                      ))}
                      <th
                        scope="col"
                        className="relative py-3.5 pl-3 pr-4 sm:pr-0"
                      >
                        <span className="sr-only">Edit</span>
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200">
                    {items.map((r) => (
                      <tr key={r.id}>
                        <td className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-0">
                          <pre>{r.id}</pre>
                        </td>
                        {columns.map((c) => (
                          <td
                            className="whitespace-nowrap px-3 py-4 text-sm text-gray-500"
                            key={c.columnKey}
                          >
                            {c.render(r)}
                          </td>
                        ))}
                        <td className="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-0 space-x-4">
                          {refActions?.(r)}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )
            ) : (
              <div>Loading ...</div>
            )}
          </div>
        </div>
      </div>
    </section>
  );
}
