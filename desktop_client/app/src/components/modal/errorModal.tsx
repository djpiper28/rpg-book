import { Modal as MantineModal } from "@mantine/core";
import { type ModalProps } from "./props";

export function ErrorModal(props: Readonly<ModalProps>) {
  return (
    <MantineModal
      centered
      color="red"
      onClose={props.close}
      opened={props.opened}
      size="md"
      title={props.title}
    >
      {props.children}
    </MantineModal>
  );
}
