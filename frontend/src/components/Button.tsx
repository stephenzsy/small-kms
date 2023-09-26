import classNames from "classnames";
import React from "react";

export type ButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: "primary" | "secondary" | "soft";
  color?: "default" | "danger";
  size?: "xs" | "sm" | "md" | "lg" | "xl";
};

export const Button = React.forwardRef(
  (props: ButtonProps, ref: React.ForwardedRef<HTMLButtonElement>) => {
    const {
      type = "button",
      className,
      variant = "secondary",
      size = "md",
      color = "default",
      ...rest
    } = props;
    return (
      <button
        ref={ref}
        type="button"
        className={classNames(
          "rounded-md text-sm font-semibold shadow-sm",
          {
            "focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2":
              variant === "primary",
            "bg-white shadow-sm ring-1 ring-inset": variant === "secondary",
          },
          {
            "bg-indigo-600 text-white hover:bg-indigo-500 focus-visible:outline-indigo-600":
            variant === "primary" && color === "default",
            "bg-red-600 text-white hover:bg-red-500 focus-visible:outline-red-600":
              variant === "primary" && color === "danger",
            "text-neutral-900 ring-neutral-300 hover:bg-neutral-50":
              variant === "secondary" && color === "default",
            "text-red-900 ring-red-300 hover:bg-red-50":
              variant === "secondary" && color === "danger",
            "bg-red-50  text-red-600  hover:bg-red-100":
              variant === "soft" && color === "danger",
          },
          {
            "px-2.5 py-1.5": size === "md",
            "px-3.5 py-2.5": size === "xl",
          },
          className
        )}
        {...rest}
      />
    );
  }
);
