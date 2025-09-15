import { Button, TextInput } from "@mantine/core";
import { type ReactNode, useEffect, useState } from "react";
import { IconSelector } from "@/components/input/iconSelector";
import { MarkdownEditor } from "@/components/input/markdownEditor";
import { read } from "@/lib/electron";
import { getProjectClient } from "@/lib/grpcClient/client";
import {
  type BasicCharacterDetails,
  type CharacterHandle,
} from "@/lib/grpcClient/pb/project_character";
import { useGlobalErrorStore } from "@/stores/globalErrorStore";
import { useTabStore } from "@/stores/tabStore";

interface Props {
  characterHandle: CharacterHandle;
  closeDialog: () => void;
}

export default function CreateCharacterModal(
  props: Readonly<Props>,
): ReactNode {
  const [characterDetails, setCharacterDetails] =
    useState<BasicCharacterDetails>({
      description: "",
      icon: new Uint8Array(),
      name: "",
    });

  const [icon, setIcon] = useState<string>("");
  const { setError } = useGlobalErrorStore((x) => x);
  const projectHandle = useTabStore((x) => x.selectedTab);

  useEffect(() => {
    if (!projectHandle) {
      return;
    }

    getProjectClient()
      .getCharacter({
        character: props.characterHandle,
        project: projectHandle,
      })
      .then((resp) => {
        setCharacterDetails(resp.response);
        const buffer = resp.response.icon;

        const binary = new Uint8Array(buffer).reduce(
          (acc, byte) => acc + String.fromCodePoint(byte),
          "",
        );

        setIcon(btoa(binary));
      })
      .catch((error: unknown) => {
        setError({
          body: String(error),
        });
      });
  }, [props.characterHandle, projectHandle, setError]);

  if (!projectHandle) {
    return "No project selected";
  }

  return (
    <div className="flex flex-col gap-3">
      <div className="flex flex-row gap-3 justify-between items-start">
        <TextInput
          label="Character Name"
          onChange={(x) => {
            const details = structuredClone(characterDetails);
            details.name = x.target.value;
            setCharacterDetails(details);
          }}
          placeholder="John Smith"
          required={true}
          value={characterDetails.name}
        />
        <IconSelector
          description="Select an icon for your character"
          filepath={icon}
          setFilepath={(src) => {
            setIcon(src);
          }}
        />
      </div>
      <MarkdownEditor
        label="Description and notes"
        setValue={(value) => {
          const details = structuredClone(characterDetails);
          details.description = value;
          setCharacterDetails(details);
        }}
        value={characterDetails.description}
      />
      <Button
        onClick={() => {
          const f = async (): Promise<void> => {
            try {
              const FileProtocol = "file://";

              if (icon.includes(FileProtocol, 0)) {
                const iconBytes = await read(icon.replace(FileProtocol, ""));
                characterDetails.icon = iconBytes;
              }

              await getProjectClient().updateCharacter({
                details: characterDetails,
                handle: props.characterHandle,
                project: projectHandle,
              });

              props.closeDialog();
            } catch (error: unknown) {
              setError({
                body: String(error),
              });
            }
          };

          // eslint-disable-next-line @typescript-eslint/no-floating-promises
          f();
        }}
      >
        Save
      </Button>
    </div>
  );
}
