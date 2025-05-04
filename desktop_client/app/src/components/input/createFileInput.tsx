import * as remote from "@electron/remote";
import { Input, type InputVariant, InputWrapper } from "@mantine/core";
import { FileIcon } from "lucide-react";

interface Props {
  buttonLabel?: string;
  description?: string;
  error?: string;
  filters: Electron.FileFilter[];
  label?: string;
  onChange: (value: string) => void;
  placeholder?: string;
  title?: string;
  value: string;
  variant?: InputVariant;
  withAsterisk?: boolean;
}

export function CreateFileInput(props: Readonly<Props>) {
  return (
    <InputWrapper
      description={props.description}
      label={props.label}
      withAsterisk={props.withAsterisk}
    >
      <Input
        error={props.error}
        onChange={(e) => {
          props.onChange(e.target.value);
        }}
        placeholder={props.placeholder}
        rightSection={
          <button
            onClick={() => {
              remote.dialog
                .showSaveDialog(remote.getCurrentWindow(), {
                  buttonLabel: props.buttonLabel,
                  filters: props.filters,
                  title: props.title,
                })
                .then((result) => {
                  if (result.canceled) {
                    return;
                  }

                  props.onChange(result.filePath);
                })
                .catch(console.error);
            }}
          >
            <FileIcon />
          </button>
        }
        rightSectionPointerEvents="all"
        value={props.value}
      />
    </InputWrapper>
  );
}
