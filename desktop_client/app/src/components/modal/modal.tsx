import { Modal as MantineModal } from "@mantine/core";
import { type ReactNode } from "react";

interface Props {
  children: ReactNode;
  // useDisclosure() from mantine/hooks to get this
  close: () => void;
  // useDisclosure() from mantine/hooks to get this
  opened: boolean;
  title: string;
}

export function Modal(props: Readonly<Props>) {
  return (
    <MantineModal
      centered
      onClose={props.close}
      opened={props.opened}
      size="auto"
      title={props.title}
    >
      {props.children}
    </MantineModal>
  );
}
