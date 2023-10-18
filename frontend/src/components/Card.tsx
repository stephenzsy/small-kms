import React, { PropsWithChildren } from "react";

type CardTitleProps = {
  description?: React.ReactNode;
};

export function Card(props: PropsWithChildren<CardTitleProps>) {
  return <div className="bg-white shadow sm:rounded-lg">{props.children}</div>;
}

export function CardTitle(props: PropsWithChildren<CardTitleProps>) {
  return (
    <div className="px-4 py-6 sm:px-6">
      <h3 className="text-base font-semibold leading-7 text-gray-900">
        {props.children}
      </h3>
      {props.description && (
        <p className="mt-1 max-w-2xl text-sm leading-6 text-gray-500">
          {props.description}
        </p>
      )}
    </div>
  );
}

export function CardSection(props: PropsWithChildren<{}>) {
  return (
    <div className="border-t border-neutral-200 px-4 py-6 max-w-full">
      {props.children}
    </div>
  );
}
