import { Input, type InputVariant, InputWrapper } from "@mantine/core";
import { FileIcon } from "lucide-react";
import { electronDialog as dialog } from "@/lib/electron";

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
              dialog
                .showSaveDialog({
                  buttonLabel: props.buttonLabel,
                  filters: props.filters,
                  title: props.title,
                })
                .then((result: Electron.SaveDialogReturnValue) => {
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
