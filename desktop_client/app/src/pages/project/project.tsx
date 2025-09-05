import { Tabs } from "@mantine/core";
import { useNavigate } from "react-router";
import { H2 } from "@/components/typography/H2";
import { mustVoid } from "@/lib/utils/errorHandlers";
import { useTabStore } from "@/stores/tabStore";
import { indexPath } from "../path";
import { CharacterTab } from "./tabs/characterTab";

export function Component() {
  const navigate = useNavigate();
  const tabs = useTabStore();

  if (!tabs.selectedTab) {
    mustVoid(navigate(indexPath));
    return;
  }

  const currentTab = tabs.tabs[tabs.selectedTab.id];

  return (
    <>
      <Tabs defaultValue="characters">
        <Tabs.List>
          <H2>{currentTab.name}</H2>
          <Tabs.Tab value="characters">Characters</Tabs.Tab>
          {/*
          <Tabs.Tab value="notes">Notes</Tabs.Tab>
          <Tabs.Tab value="items">Items</Tabs.Tab>
          <Tabs.Tab value="factions">Factions</Tabs.Tab>
          <Tabs.Tab value="maps">Maps &amp; Places</Tabs.Tab>
          <Tabs.Tab value="settings">Settings</Tabs.Tab>
          */}
        </Tabs.List>

        <Tabs.Panel value="characters">
          <CharacterTab />
        </Tabs.Panel>
      </Tabs>
    </>
  );
}
