import { Modal as MantineModal } from "@mantine/core";
import { type ModalProps } from "./props";

export function Modal(props: Readonly<ModalProps>) {
  return (
    <MantineModal
      centered
      onClose={props.close}
      opened={props.opened}
      size="xl"
      title={props.title}
    >
      {props.children}
    </MantineModal>
  );
}
