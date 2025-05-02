import { type ReactNode } from "react";

export interface ModalProps {
  children: ReactNode;
  // useDisclosure() from mantine/hooks to get this
  close: () => void;
  // useDisclosure() from mantine/hooks to get this
  opened: boolean;
  title: string;
}
