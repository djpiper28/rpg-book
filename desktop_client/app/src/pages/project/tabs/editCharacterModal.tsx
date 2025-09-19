import { Button } from "@mantine/core";
import { type ReactNode, useEffect, useState } from "react";
import { getProjectClient } from "@/lib/grpcClient/client";
import {
  type BasicCharacterDetails,
  type CharacterHandle,
} from "@/lib/grpcClient/pb/project_character";
import { base64ToUint8Array, uint8ArrayToBase64 } from "@/lib/utils/base64";
import { useGlobalErrorStore } from "@/stores/globalErrorStore";
import { useTabStore } from "@/stores/tabStore";
import { CharacterEdit } from "./characterEdit";

interface Props {
  characterHandle: CharacterHandle;
  closeDialog: () => void;
}

export default function EditCharacterModal(props: Readonly<Props>): ReactNode {
  const [characterDetails, setCharacterDetails] =
    useState<BasicCharacterDetails>({
      description: "",
      icon: new Uint8Array(),
      name: "",
    });

  const [iconB64, setIconB64] = useState<string>("");
  const [dirtyIcon, setDirtyIcon] = useState(false);
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

        if (buffer.length > 0) {
          const b64 = uint8ArrayToBase64(buffer);
          setIconB64(`data:image/png;base64,${b64}`);
          setDirtyIcon(true);
        }
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
      <CharacterEdit
        character={characterDetails}
        iconB64={iconB64}
        setCharacter={setCharacterDetails}
        setIconB64={setIconB64}
      />
      <Button
        onClick={() => {
          const f = async (): Promise<void> => {
            try {
              if (dirtyIcon) {
                const b64 = iconB64.split(",")[1];
                characterDetails.icon = base64ToUint8Array(b64);
              }

              await getProjectClient().updateCharacter({
                details: characterDetails,
                handle: props.characterHandle,
                project: projectHandle,
                setImage: dirtyIcon,
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
