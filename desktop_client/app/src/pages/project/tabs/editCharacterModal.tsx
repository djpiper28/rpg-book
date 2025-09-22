import { Button } from "@mantine/core";
import { type ReactNode, useEffect, useState } from "react";
import { getLogger, getProjectClient } from "@/lib/grpcClient/client";
import {
  type BasicCharacterDetails,
  type CharacterHandle,
} from "@/lib/grpcClient/pb/project_character";
import { uint8ArrayToBase64 } from "@/lib/utils/base64";
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
      handle: {
        id: "",
      },
      icon: new Uint8Array(),
      iconPath: "",
      name: "",
    });

  const [iconPath, setIconPath] = useState<string>("");
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

  const iconB64 = uint8ArrayToBase64(characterDetails.icon);

  return (
    <div className="flex flex-col gap-3">
      <CharacterEdit
        character={characterDetails}
        iconPath={iconPath}
        imageDataB64={iconB64}
        setCharacter={setCharacterDetails}
        setIconPath={(path) => {
          setIconPath(path);
          setDirtyIcon(true);
        }}
      />
      <Button
        onClick={() => {
          const f = async (): Promise<void> => {
            try {
              await getProjectClient().updateCharacter({
                details: {
                  ...characterDetails,
                  iconPath,
                },
                handle: props.characterHandle,
                project: projectHandle,
                setImage: dirtyIcon,
              });

              props.closeDialog();
            } catch (error: unknown) {
              getLogger().error("Cannot edit character", {
                error: String(error),
              });

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
