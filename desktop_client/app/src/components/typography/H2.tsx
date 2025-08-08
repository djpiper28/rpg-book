import { type ReactNode } from "react";

interface Props {
  children?: string | ReactNode;
}

export function H2(props: Readonly<Props>) {
  return <h2 className="text-2xl font-semibold">{props.children}</h2>;
}
