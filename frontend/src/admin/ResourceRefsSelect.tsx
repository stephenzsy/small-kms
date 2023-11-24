import Select, { SelectProps } from "antd/es/select";
import { useMemo } from "react";
import { Ref as ResourceReference } from "../generated/apiv2";

export type ResourceRefsSelectProps = Omit<
  SelectProps<string>,
  "options"
> & {
  data: ResourceReference[] | undefined;
};

export function ResourceRefsSelect({
  data,
  ...restProps
}: ResourceRefsSelectProps) {
  const options = useMemo(() => {
    return data?.map((profile) => {
      return {
        label: (
          <span>
            {profile.displayName} ({profile.id})
          </span>
        ),
        value: profile.id,
        profile,
      };
    });
  }, [data]);

  return (
    <>
      <Select options={options} {...restProps} />
    </>
  );
}
