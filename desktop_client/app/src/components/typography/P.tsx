import { type ReactNode } from "react";

interface Props {
  children: string | ReactNode;
  className?: string;
}

export function P(props: Readonly<Props>) {
  return <p className="text-lg ${props.className}">{props.children}</p>;
}
