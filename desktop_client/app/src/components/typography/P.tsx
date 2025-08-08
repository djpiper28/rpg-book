import { type ReactNode } from "react";

interface Props {
  children?: string | ReactNode;
  className?: string;
}

export function P(props: Readonly<Props>) {
  // eslint-disable-next-line @typescript-eslint/no-unnecessary-template-expression
  return <p className={`${props.className ?? ""}`}>{props.children}</p>;
}
