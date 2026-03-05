import { type ReactNode, useEffect, useState } from "react";
import { EditDelete } from "@/components/buttons/editDelete";
import MarkdownRenderer from "@/components/renderers/markdown";
import { H2 } from "@/components/typography/H2";
import { getLogger, getProjectClient } from "@/lib/grpcClient/client";
import {
  type CharacterDetails,
  type CharacterHandle,
} from "@/lib/grpcClient/pb/project_character";
import { useTabStore } from "@/stores/tabStore";

export interface CharacterViewProps {
  characterHandle: CharacterHandle;
  isEditModalOpen: boolean;
  onDelete: () => void;
  onEdit: () => void;
}

export function CharacterView(props: Readonly<CharacterViewProps>): ReactNode {
  const [character, setCharacter] = useState<CharacterDetails | undefined>();
  const [iconUrl, setIconUrl] = useState("");
  const projectHandle = useTabStore((x) => x.selectedTab);

  useEffect(() => {
    // eslint-disable-next-line @typescript-eslint/explicit-function-return-type
    const sync = async () => {
      if (!projectHandle) {
        return;
      }

      try {
        const resp = await getProjectClient().getCharacter({
          character: props.characterHandle,
          project: projectHandle,
        });

        setCharacter(resp.response.details);
        setIconUrl(resp.response.iconUrl);
      } catch (error: unknown) {
        getLogger().error("Cannot get character", {
          character: props.characterHandle.id,
          error: JSON.stringify(error),
          project: JSON.stringify(projectHandle),
        });
      }
    };

    // eslint-disable-next-line @typescript-eslint/no-floating-promises
    sync();
  }, [projectHandle, props.characterHandle, props.isEditModalOpen]);

  if (!character) {
    return <div className="flex-1" />;
  }

  return (
    <div className="flex-1 gap-3 overflow-y-auto flex flex-col">
      <div className="flex flex-col gap-3">
        <div className="flex flex-row gap-3 items-center justify-between">
          <H2>{character.name}</H2>
          <EditDelete delete={props.onDelete} edit={props.onEdit} />
        </div>
        {iconUrl && (
          <img alt="User selected" className="max-h-screen" src={iconUrl} />
        )}
      </div>
      <MarkdownRenderer
        className="flex-grow min-h-0"
        markdown={character.description}
      />
    </div>
  );
}
