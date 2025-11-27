import { type ReactNode } from "react";

interface Props {
  children?: string | ReactNode;
  id?: string;
}

export function H2(props: Readonly<Props>): ReactNode {
  return (
    <h2 className="text-2xl font-semibold" id={props.id}>
      {props.children}
    </h2>
  );
}
