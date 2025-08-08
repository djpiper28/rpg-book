import { shell } from "@electron/remote";
import { type ReactNode } from "react";
import { logger } from "@/lib/grpcClient/client";

interface Props {
  children: ReactNode;
  href: string;
  newWindow?: boolean;
  openInBrowser?: boolean;
}

export function Link(props: Readonly<Props>) {
  const style = "text-blue-500 hover:text-blue-600 underline";

  return props.openInBrowser ? (
    <button
      className={`${style} cursor-grab`}
      onClick={() => {
        shell
          .openExternal(props.href)
          .then()
          .catch((error: unknown) => {
            // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
            logger.error("Cannot open link", { error: `${error}` });
          });
      }}
    >
      {props.children}
    </button>
  ) : (
    <a
      className={style}
      href={props.href}
      rel={props.newWindow ? "noreferrer" : ""}
      target={props.newWindow ? "_blank" : ""}
    >
      {props.children}
    </a>
  );
}
