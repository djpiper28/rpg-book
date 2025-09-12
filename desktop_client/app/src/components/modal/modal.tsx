import { Modal as MantineModal } from "@mantine/core";
import { type ReactNode } from "react";
import { type ModalProps } from "./props";

export function Modal(props: Readonly<ModalProps>): ReactNode {
  return (
    <MantineModal
      centered
      onClose={props.close}
      opened={props.opened}
      size={props.size ?? "xl"}
      title={props.title}
    >
      <div className="flex flex-col gap-2">{props.children}</div>
    </MantineModal>
  );
}
