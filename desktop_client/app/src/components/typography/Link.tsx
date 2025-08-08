import { type ReactNode } from "react";

interface Props {
  children: ReactNode;
  href: string;
  newWindow?: boolean;
}

export function Link(props: Readonly<Props>) {
  return (
    <a
      className="text-blue-700 hover:text-blue-600 underline"
      href={props.href}
      rel={props.newWindow ? "noreferrer" : ""}
      target={props.newWindow ? "_blank" : ""}
    >
      {props.children}
    </a>
  );
}
