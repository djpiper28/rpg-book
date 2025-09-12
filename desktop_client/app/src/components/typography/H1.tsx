import { type ReactNode } from "react";

interface Props {
  children?: string | ReactNode;
}

export function H1(props: Readonly<Props>): ReactNode {
  return <h1 className="text-3xl font-bold">{props.children}</h1>;
}
