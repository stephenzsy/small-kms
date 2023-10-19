import { Typography } from "antd";
import classNames from "classnames";
import { LinkProps, Link as RRLink } from "react-router-dom";

export function Link({ className, ...restProps }: LinkProps) {
  return (
    <RRLink
      className={classNames("text-indigo-600 hover:text-indigo-900", className)}
      {...restProps}
    />
  );
}
