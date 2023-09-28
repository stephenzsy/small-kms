import { type SetStateAction, useState } from "react";

export type ValueState<T> = {
  value: T;
  onChange: React.Dispatch<SetStateAction<T>>;
};

export function useValueState<T = string>(
  initialValue: T | (() => T)
): ValueState<T> {
  const [state, setState] = useState(initialValue);
  return {
    value: state,
    onChange: setState,
  };
}

export type FixedValueState<T> = {
  value: T;
  onChange: undefined;
};

export type ValueStateMayBeFixed<T> = ValueState<T> | FixedValueState<T>;

export function useFixedValueState<T>(
  valueState: ValueState<T>,
  fixedValue: T | undefined
): ValueStateMayBeFixed<T> {
  if (fixedValue === undefined) {
    return valueState;
  }
  return { value: fixedValue, onChange: undefined };
}
