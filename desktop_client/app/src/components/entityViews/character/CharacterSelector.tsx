import { Checkbox } from "@mantine/core";
import { type ReactNode, useEffect, useState } from "react";
import { Search } from "@/components/search/search";
import { H3 } from "@/components/typography/H3";
import { getLogger, getProjectClient } from "@/lib/grpcClient/client";
import {
  type CharacterDetails,
  type CharacterHandle,
} from "@/lib/grpcClient/pb/project_character";
import { useProjectStore } from "@/stores/projectStore";
import { useTabStore } from "@/stores/tabStore";

interface Props {
  onSelectedCharactersChanges: (handles: CharacterHandle[]) => void;
  selectedCharacters: CharacterHandle[];
}

export function CharacterSelector(props: Readonly<Props>): ReactNode {
  const projectHandle = useTabStore((x) => x.selectedTab);
  const projectStore = useProjectStore((x) => x);
  const thisProject = projectHandle && projectStore.getProject(projectHandle);
  const [queryError, setQueryError] = useState("");
  const [queryText, setQueryText] = useState("");
  const [queryResult, setQueryResult] = useState<CharacterDetails[]>([]);

  useEffect(() => {
    // eslint-disable-next-line @typescript-eslint/explicit-function-return-type
    const sync = async () => {
      if (!projectHandle) {
        return;
      }

      if (!thisProject) {
        return;
      }

      if (queryText.trim() === "") {
        setQueryResult(thisProject.project.characters);
        setQueryError("");
        return;
      }

      try {
        const resp = await getProjectClient().searchCharacter({
          project: projectHandle,
          query: queryText,
        });

        setQueryResult(
          thisProject.project.characters.filter((allCharacters) => {
            const inSearchRes = resp.response.details.find(
              (c) => c.id === allCharacters.handle?.id,
            );

            return !!inSearchRes;
          }),
        );

        setQueryError("");
      } catch (error: unknown) {
        getLogger().error("Cannot search characters", {
          error: JSON.stringify(error),
          project: JSON.stringify(projectHandle),
          query: queryText,
        });

        setQueryError(JSON.stringify(error));
      }
    };

    // eslint-disable-next-line @typescript-eslint/no-floating-promises
    sync();
  }, [projectHandle, queryText, thisProject, thisProject?.project.characters]);

  return (
    <div className="h-1/3">
      <H3>Character Selector</H3>
      <div className="flex flex-col border-2 rounded-2xl border-gray-500 p-2">
        <Search<CharacterDetails>
          elementWrapper={(children: ReactNode): ReactNode => {
            return <div className="flex flex-col gap-2">{children}</div>;
          }}
          error={queryError}
          onChange={(txt: string) => {
            setQueryText(txt);
          }}
          placeholder="Character search"
          render={(character: CharacterDetails): ReactNode => {
            return (
              <Checkbox
                checked={props.selectedCharacters.some(
                  (x) => x.id === character.handle?.id,
                )}
                label={character.name}
                onChange={(ev) => {
                  const handle = character.handle;

                  if (!handle) {
                    return;
                  }

                  if (ev.target.checked) {
                    const newSelected = structuredClone(
                      props.selectedCharacters,
                    );

                    if (!newSelected.some((x) => x.id === handle.id)) {
                      newSelected.push(handle);
                      props.onSelectedCharactersChanges(newSelected);
                    }
                  } else {
                    const newSelected = structuredClone(
                      props.selectedCharacters,
                    ).filter((x) => x.id === character.handle?.id);

                    props.onSelectedCharactersChanges(newSelected);
                  }
                }}
              />
            );
          }}
          searchRes={queryResult}
        />
      </div>
    </div>
  );
}
