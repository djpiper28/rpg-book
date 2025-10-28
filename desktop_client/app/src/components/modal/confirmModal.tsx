import { Button, Modal as MantineModal } from "@mantine/core";
import { type ReactNode } from "react";

interface Props {
  // useDisclosure() from mantine/hooks to get this
  close: () => void;
  onConfirm: () => void;
  // useDisclosure() from mantine/hooks to get this
  opened: boolean;
  title: string;
}

export function ConfirmModal(props: Readonly<Props>): ReactNode {
  return (
    <MantineModal
      centered
      onClose={props.close}
      opened={props.opened}
      size="sm"
      title={props.title}
    >
      <div className="flex flex-row gap-2 justify-end-safe">
        <Button
          color="red"
          onClick={() => {
            props.close();
          }}
        >
          No
        </Button>
        <Button
          color="blue"
          onClick={() => {
            props.onConfirm();
            props.close();
          }}
        >
          Yes
        </Button>
      </div>
    </MantineModal>
  );
}
