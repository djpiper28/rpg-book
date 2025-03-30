import { ReactNode } from "react";
import React from "react";

interface Props {
  href: string;
  children: ReactNode;
}

export function Link(props: Readonly<Props>) {
  return (
    <a
      href={props.href}
      className="text-blue-700 hover:text-blue-600 underline"
    >
      {props.children}
    </a>
  );
}
