import { Center, Loader, Tabs } from "@mantine/core";
import { type ReactNode, Suspense, lazy } from "react";
import { useNavigate } from "react-router";
import { H2 } from "@/components/typography/H2";
import { mustVoid } from "@/lib/utils/errorHandlers";
import { useTabStore } from "@/stores/tabStore";
import { indexPath } from "../path";

const CharacterTab = lazy(() =>
  import("./tabs/characters/characterTab").then((m) => ({
    default: m.CharacterTab,
  })),
);

const NoteTab = lazy(() =>
  import("./tabs/notes/noteTab").then((m) => ({ default: m.NoteTab })),
);

const defaultValue = "characters";

export function Component(): ReactNode {
  const navigate = useNavigate();
  const tabs = useTabStore();

  if (!tabs.selectedTab) {
    mustVoid(navigate(indexPath));
    return;
  }

  const currentTab = tabs.tabs[tabs.selectedTab.id];

  return (
    <>
      <Tabs
        defaultValue="characters"
        keepMounted={false}
        onChange={(value: string | null) => {
          value ??= defaultValue;

          tabs.setValue(tabs.selectedTab?.id ?? "", value);
        }}
        value={currentTab.value}
      >
        <Tabs.List>
          <H2>{currentTab.name}</H2>
          <Tabs.Tab value={defaultValue}>Characters</Tabs.Tab>
          <Tabs.Tab value="notes">Notes</Tabs.Tab>
          {/*
          <Tabs.Tab value="items">Items</Tabs.Tab>
          <Tabs.Tab value="factions">Factions</Tabs.Tab>
          <Tabs.Tab value="maps">Maps &amp; Places</Tabs.Tab>
          <Tabs.Tab value="settings">Settings</Tabs.Tab>
          */}
        </Tabs.List>

        <Tabs.Panel value="characters">
          <Suspense
            fallback={
              <Center p="md">
                <Loader />
              </Center>
            }
          >
            <CharacterTab />
          </Suspense>
        </Tabs.Panel>

        <Tabs.Panel value="notes">
          <Suspense
            fallback={
              <Center p="md">
                <Loader />
              </Center>
            }
          >
            <NoteTab />
          </Suspense>
        </Tabs.Panel>
      </Tabs>
    </>
  );
}
