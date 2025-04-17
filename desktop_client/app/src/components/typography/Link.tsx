import { type ReactNode } from "react";

interface Props {
  children: ReactNode;
  href: string;
}

export function Link(props: Readonly<Props>) {
  return (
    <a
      className="text-blue-700 hover:text-blue-600 underline"
      href={props.href}
    >
      {props.children}
    </a>
  );
}
