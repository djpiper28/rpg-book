import { Modal as MantineModal } from "@mantine/core";
import { type ReactNode } from "react";
import { type ModalProps } from "./props";

export function ErrorModal(props: Readonly<ModalProps>): ReactNode {
  return (
    <MantineModal
      centered
      onClose={props.close}
      opened={props.opened}
      size="md"
      title={props.title}
    >
      {props.children}
    </MantineModal>
  );
}
