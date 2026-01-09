import { type ReactNode } from "react";

interface Props {
  children?: string | ReactNode;
  id?: string;
}

export function H3(props: Readonly<Props>): ReactNode {
  return (
    <h3 className="text-xl font-semibold" id={props.id}>
      {props.children}
    </h3>
  );
}
