import { Button } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { TriangleAlert } from "lucide-react";
import { type ReactNode } from "react";
import { electronShell } from "@/lib/electron";
import { getLogger } from "@/lib/grpcClient/client";
import { Modal } from "../modal/modal";
import { P } from "./P";

interface Props {
  children: ReactNode;
  href: string;
  // Opens in a new electron window
  newWindow?: boolean;
  openInBrowser?: boolean;
  // Set if using openInBrowser and the URL is known safe
  safe?: boolean;
}

export function Link(props: Readonly<Props>) {
  const [opened, { close, open }] = useDisclosure(false);
  const style = "text-blue-500 hover:text-blue-600 underline";

  if (!props.openInBrowser) {
    return (
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

  let url: URL | undefined;

  try {
    url = new URL(props.href);
  } catch (error: unknown) {
    return <p className={style}>Invalid URL</p>;
  }

  const openUrl = () => {
    electronShell
      .openExternal(props.href)
      .then()
      .catch((error: unknown) => {
        // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
        getLogger().error("Cannot open link", { error: `${error}` });
      });
  };

  return (
    <>
      <Modal
        close={close}
        opened={opened}
        size="auto"
        title="Opening External Links Can Be Dangerous"
      >
        <div className="flex flex-row gap-2">
          <TriangleAlert />
          <P>
            You are about to open an external link, this can be dangerous. Make
            sure you are happy with the website before opening it. ({url.href})
          </P>
        </div>
        <Button
          onClick={() => {
            openUrl();
          }}
        >
          Open Anyway
        </Button>
      </Modal>
      <button
        className={`${style} cursor-grab`}
        onClick={() => {
          if (props.safe) {
            openUrl();
          } else {
            open();
          }
        }}
      >
        {props.children}
      </button>
    </>
  );
}
